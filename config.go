package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var flagHelpError = flag.ErrHelp

type Config struct {
	Root                  string   `json:"root"`
	ExtraRoots            []string `json:"extra_roots"`
	Addr                  string   `json:"addr"`
	ThumbDir              string   `json:"thumb_dir"`
	Meta                  string   `json:"meta"`
	FFmpeg                string   `json:"ffmpeg"`
	ThumbSeek             float64  `json:"thumb_seek"`
	ThumbWidth            int      `json:"thumb_width"`
	ThumbWorkers          int      `json:"thumb_workers"`
	ThumbMaintenanceEvery string   `json:"thumb_maintenance_every"`
	SizeRefresh           string   `json:"size_refresh"`
	AutoDeleteDays        int      `json:"auto_delete_days"`
	AutoDeleteEvery       string   `json:"auto_delete_every"`
	DetectionEnabled      bool     `json:"detection_enabled"`
	DetectionPython       string   `json:"detection_python"`
	DetectionModel        string   `json:"detection_model"`
	DetectionThreads      int      `json:"detection_threads"`
	DetectionFPS          float64  `json:"detection_fps"`
	DetectionConfidence   float64  `json:"detection_confidence"`
	DetectionMinFrames    int      `json:"detection_min_frames"`
	DetectionMergeGap     float64  `json:"detection_merge_gap"`
	DetectionCarMovement  float64  `json:"detection_car_movement"`
	DetectionEvery        string   `json:"detection_every"`
	DetectionMinAge       string   `json:"detection_min_age"`
	DetectionTimeout      string   `json:"detection_timeout"`
	DetectionPrune        string   `json:"detection_prune"`
}

func defaultConfig() Config {
	return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h",
		DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
}

func decodeConfig(data []byte) (Config, error) {
	c := defaultConfig()
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return c, errors.New("config must be a JSON object")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return c, err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return c, errors.New("config must contain exactly one JSON object")
	}
	if c.ExtraRoots == nil {
		c.ExtraRoots = []string{}
	}
	return c, nil
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Root) == "" {
		return errors.New("root cannot be empty")
	}
	for _, p := range append([]string{c.Root}, c.ExtraRoots...) {
		if !filepath.IsAbs(expandPath(p)) {
			return fmt.Errorf("recording directory must be absolute: %s", p)
		}
		st, err := os.Stat(expandPath(p))
		if err != nil {
			return fmt.Errorf("recording directory %s: %w", p, err)
		}
		if !st.IsDir() {
			return fmt.Errorf("recording directory %s is not a directory", p)
		}
	}
	for name, p := range map[string]string{"thumb_dir": c.ThumbDir, "meta": c.Meta} {
		if p != "" && !filepath.IsAbs(expandPath(p)) {
			return fmt.Errorf("%s must be absolute or empty", name)
		}
	}
	host, port, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return fmt.Errorf("addr: %w", err)
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return errors.New("addr must use a port between 1 and 65535")
	}
	if _, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port)); err != nil {
		return fmt.Errorf("addr: %w", err)
	}
	if strings.TrimSpace(c.FFmpeg) == "" {
		return errors.New("ffmpeg cannot be empty")
	}
	if c.ThumbWorkers < 1 || c.ThumbWorkers > 64 {
		return errors.New("thumb_workers must be between 1 and 64")
	}
	if c.ThumbWidth < 64 || c.ThumbWidth > 4096 {
		return errors.New("thumb_width must be between 64 and 4096")
	}
	if math.IsNaN(c.ThumbSeek) || math.IsInf(c.ThumbSeek, 0) || c.ThumbSeek < 0 || c.ThumbSeek > 3600 {
		return errors.New("thumb_seek must be between 0 and 3600 seconds")
	}
	if c.AutoDeleteDays < 0 || c.AutoDeleteDays > 36500 {
		return errors.New("auto_delete_days must be between 0 and 36500; 0 disables deletion")
	}
	if err := c.validateDetection(); err != nil {
		return err
	}
	for name, v := range map[string]string{"thumb_maintenance_every": c.ThumbMaintenanceEvery, "size_refresh": c.SizeRefresh, "auto_delete_every": c.AutoDeleteEvery} {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("%s: enter a duration such as 30m or 6h", name)
		}
		min := time.Minute
		if name == "size_refresh" {
			min = 10 * time.Second
		}
		if name == "thumb_maintenance_every" && d == 0 {
			continue
		}
		if d < min {
			return fmt.Errorf("%s must be at least %s", name, min)
		}
	}
	meta := c.Meta
	if meta == "" {
		meta = filepath.Join(expandPath(c.Root), ".blinkclips_meta.json")
	}
	if _, err := NewMetaStore(expandPath(meta)); err != nil {
		return fmt.Errorf("favorites metadata: %w", err)
	}
	return nil
}

func configRevision(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }

// The entire parent directory is mounted into containers so rename remains atomic.
func atomicWriteConfig(path string, data []byte) error {
	if st, err := os.Lstat(path); err == nil && !st.Mode().IsRegular() {
		return errors.New("config path must be a regular file")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	original, _ := os.Stat(path)
	if original == nil && strings.HasSuffix(path, ".previous") {
		original, _ = os.Stat(strings.TrimSuffix(path, ".previous"))
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".neonvr-config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if original != nil && os.Geteuid() == 0 {
		if st, ok := original.Sys().(*syscall.Stat_t); ok {
			if err = f.Chown(int(st.Uid), int(st.Gid)); err != nil {
				f.Close()
				return err
			}
		}
	}
	if err = f.Chmod(0640); err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err = os.Rename(f.Name(), path); err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

func configBytes(c Config) []byte {
	data, _ := json.MarshalIndent(c, "", "  ")
	return append(data, '\n')
}

func loadStartupConfig(args []string) (Config, string, error) {
	path := "config/neonvr.json"
	explicit := false
	for i, a := range args {
		if a == "-config" || a == "--config" {
			if i+1 >= len(args) {
				return Config{}, "", errors.New("-config requires a path")
			}
			path = args[i+1]
			explicit = true
		}
		if strings.HasPrefix(a, "-config=") || strings.HasPrefix(a, "--config=") {
			path = strings.SplitN(a, "=", 2)[1]
			explicit = true
		}
	}
	if strings.TrimSpace(path) == "" {
		return Config{}, "", errors.New("config path cannot be empty")
	}
	path, err := filepath.Abs(expandPath(path))
	if err != nil {
		return Config{}, "", err
	}
	c := defaultConfig()
	raw, readErr := os.ReadFile(path)
	if readErr == nil {
		c, err = decodeConfig(raw)
		if err != nil {
			return c, path, fmt.Errorf("config %s: %w", path, err)
		}
	} else if !os.IsNotExist(readErr) || explicit {
		return c, path, readErr
	}
	fs := flag.NewFlagSet("neonvr", flag.ContinueOnError)
	fs.StringVar(&path, "config", path, "JSON config file (default config/neonvr.json)")
	fs.StringVar(&c.Root, "root", c.Root, "primary recordings directory")
	extra := stringListFlag(append([]string{}, c.ExtraRoots...))
	hasExtra := false
	for _, a := range args {
		if a == "-extra-root" || a == "--extra-root" || strings.HasPrefix(a, "-extra-root=") || strings.HasPrefix(a, "--extra-root=") {
			hasExtra = true
		}
	}
	if hasExtra {
		extra = nil
	}
	fs.Var(&extra, "extra-root", "additional directory; repeat to replace configured extra roots")
	fs.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	fs.StringVar(&c.ThumbDir, "thumb-dir", c.ThumbDir, "thumbnail directory; empty uses <root>/.blinkclips_thumbs")
	fs.StringVar(&c.Meta, "meta", c.Meta, "favorites JSON; empty uses <root>/.blinkclips_meta.json")
	fs.StringVar(&c.FFmpeg, "ffmpeg", c.FFmpeg, "ffmpeg executable")
	fs.Float64Var(&c.ThumbSeek, "thumb-seek", c.ThumbSeek, "fallback thumbnail seek in seconds")
	fs.IntVar(&c.ThumbWidth, "thumb-width", c.ThumbWidth, "thumbnail width")
	fs.IntVar(&c.ThumbWorkers, "thumb-workers", c.ThumbWorkers, "thumbnail workers")
	fs.StringVar(&c.ThumbMaintenanceEvery, "thumb-maintenance-every", c.ThumbMaintenanceEvery, "maintenance cadence; 0 disables")
	fs.StringVar(&c.SizeRefresh, "size-refresh", c.SizeRefresh, "size statistics cache cadence")
	fs.IntVar(&c.AutoDeleteDays, "auto-delete-days", c.AutoDeleteDays, "retention days; 0 disables deletion")
	fs.StringVar(&c.AutoDeleteEvery, "auto-delete-every", c.AutoDeleteEvery, "retention check cadence")
	fs.BoolVar(&c.DetectionEnabled, "detection-enabled", c.DetectionEnabled, "enable background YOLO26n CPU detection")
	fs.StringVar(&c.DetectionPython, "detection-python", c.DetectionPython, "Python executable with ONNX Runtime installed")
	fs.StringVar(&c.DetectionModel, "detection-model", c.DetectionModel, "YOLO26n end-to-end ONNX model path")
	fs.IntVar(&c.DetectionThreads, "detection-threads", c.DetectionThreads, "CPU inference threads")
	fs.Float64Var(&c.DetectionFPS, "detection-fps", c.DetectionFPS, "sampled frames per second")
	fs.Float64Var(&c.DetectionConfidence, "detection-confidence", c.DetectionConfidence, "minimum confidence for confirmed detections")
	fs.IntVar(&c.DetectionMinFrames, "detection-min-frames", c.DetectionMinFrames, "consecutive observations required to confirm a detection")
	fs.Float64Var(&c.DetectionMergeGap, "detection-merge-gap", c.DetectionMergeGap, "maximum gap between confirmed intervals in seconds")
	fs.Float64Var(&c.DetectionCarMovement, "detection-car-movement", c.DetectionCarMovement, "car movement threshold as fraction of frame")
	fs.StringVar(&c.DetectionEvery, "detection-every", c.DetectionEvery, "recording scan interval")
	fs.StringVar(&c.DetectionMinAge, "detection-min-age", c.DetectionMinAge, "minimum time since last recording write")
	fs.StringVar(&c.DetectionTimeout, "detection-timeout", c.DetectionTimeout, "maximum analysis time per recording")
	fs.StringVar(&c.DetectionPrune, "detection-prune", c.DetectionPrune, "no-target clip policy: off, quarantine, or delete")
	if err = fs.Parse(args); err != nil {
		return c, path, err
	}
	if fs.NArg() != 0 {
		return c, path, errors.New("unexpected positional arguments")
	}
	c.ExtraRoots = append([]string{}, extra...)
	path, err = filepath.Abs(expandPath(path))
	if err != nil {
		return c, path, err
	}
	if err = c.validate(); err != nil {
		return c, path, err
	}
	// Persist explicit legacy flag overrides, so later browser restarts use those settings.
	changed := os.IsNotExist(readErr)
	if readErr == nil {
		original, _ := decodeConfig(raw)
		changed = !reflect.DeepEqual(c, original)
	}
	if changed {
		if err = atomicWriteConfig(path, configBytes(c)); err != nil {
			return c, path, err
		}
	}
	return c, path, nil
}

type configManager struct {
	mu         sync.Mutex
	path       string
	active     Config
	token      string
	restarting bool
	restart    chan struct{}
}

func newConfigManager(path string, c Config) (*configManager, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}
	return &configManager{path: path, active: c, token: hex.EncodeToString(token), restart: make(chan struct{}, 1)}, nil
}

func (m *configManager) response(w http.ResponseWriter, status int, saved Config, revision string, nextURL string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"config": saved, "active_config": m.active, "revision": revision, "token": m.token, "instance_id": m.token, "restart_required": !reflect.DeepEqual(saved, m.active), "restarting": m.restarting, "next_url": nextURL})
}

func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != "GET" && r.Method != "POST" {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", 405)
		return
	}
	if r.Method == "POST" {
		if r.Header.Get("X-NeoNVR-Token") != m.token {
			http.Error(w, "reload Settings before saving", 403)
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" {
			u, err := url.Parse(origin)
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			if err != nil || u.Host != r.Host || u.Scheme != scheme {
				http.Error(w, "same-origin requests required", 403)
				return
			}
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "JSON required", 415)
			return
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	raw, err := os.ReadFile(m.path)
	if err != nil {
		http.Error(w, "cannot read config: "+err.Error(), 500)
		return
	}
	saved, err := decodeConfig(raw)
	if err != nil {
		http.Error(w, "invalid config file: "+err.Error(), 500)
		return
	}
	revision := configRevision(raw)
	if r.Method == "GET" {
		m.response(w, 200, saved, revision, "")
		return
	}
	if m.restarting {
		http.Error(w, "restart already in progress", 409)
		return
	}
	var request struct {
		Config   json.RawMessage `json:"config"`
		Revision string          `json:"revision"`
		Restart  bool            `json:"restart"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&request); err != nil {
		http.Error(w, "invalid request: "+err.Error(), 400)
		return
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		http.Error(w, "unexpected trailing JSON", 400)
		return
	}
	if request.Revision != revision {
		http.Error(w, "settings changed elsewhere; reload before saving", 409)
		return
	}
	c, err := decodeConfig(request.Config)
	if err != nil {
		http.Error(w, "invalid config: "+err.Error(), 400)
		return
	}
	if err = c.validate(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if request.Restart && c.Addr != m.active.Addr {
		listener, err := net.Listen("tcp", c.Addr)
		if err != nil {
			http.Error(w, "cannot listen on new address: "+err.Error(), 400)
			return
		}
		listener.Close()
	}
	data := configBytes(c)
	if err = atomicWriteConfig(m.path+".previous", raw); err != nil {
		http.Error(w, "cannot save previous config: "+err.Error(), 500)
		return
	}
	if err = atomicWriteConfig(m.path, data); err != nil {
		http.Error(w, "cannot save config: "+err.Error(), 500)
		return
	}
	nextURL := ""
	if request.Restart {
		m.restarting = true
		host, port, _ := net.SplitHostPort(c.Addr)
		if host == "" || host == "0.0.0.0" || host == "::" {
			host, _, _ = net.SplitHostPort(r.Host)
			if host == "" {
				host = r.Host
			}
		}
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		nextURL = scheme + "://" + net.JoinHostPort(host, port) + "/settings"
	}
	m.response(w, 200, c, configRevision(data), nextURL)
	if request.Restart {
		// Delay until this response has reached the browser, then drain and re-exec.
		time.AfterFunc(750*time.Millisecond, func() { m.restart <- struct{}{} })
	}
}
