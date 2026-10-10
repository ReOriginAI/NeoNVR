package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const detectionMetadataTag = "neonvr_detections"
const detectionVersion = 1

var errDetectionWorker = errors.New("detection worker protocol failed")

//go:embed detection/worker.py
var detectionWorkerScript string

//go:embed static/detection-viewer.js
var detectionViewerJS []byte

type DetectionInterval struct {
	Kind       string   `json:"kind"`
	Labels     []string `json:"labels"`
	Start      float64  `json:"start"`
	End        float64  `json:"end"`
	Confidence float64  `json:"confidence"`
}

type DetectionResult struct {
	Version     int                 `json:"version"`
	Model       string              `json:"model"`
	ModelSHA256 string              `json:"model_sha256"`
	Policy      string              `json:"policy"`
	Complete    bool                `json:"complete"`
	Duration    float64             `json:"duration"`
	SampleFPS   float64             `json:"sample_fps"`
	Frames      int                 `json:"frames"`
	Uncertain   bool                `json:"uncertain"`
	Intervals   []DetectionInterval `json:"intervals"`
	AnalyzedAt  string              `json:"analyzed_at"`
}

type detectionResponse struct {
	Status     string           `json:"status"`
	Detections *DetectionResult `json:"detections,omitempty"`
	Message    string           `json:"message,omitempty"`
}

type detectionCache struct {
	Size    int64           `json:"size"`
	MtimeNS int64           `json:"mtime_ns"`
	Result  DetectionResult `json:"result"`
}

func (c Config) validateDetection() error {
	if strings.TrimSpace(c.DetectionPython) == "" || !strings.HasSuffix(strings.ToLower(c.DetectionModel), ".onnx") {
		return errors.New("detection requires a Python executable and a YOLO26n .onnx model path")
	}
	if c.DetectionThreads < 1 || c.DetectionThreads > 16 || c.DetectionMinFrames < 1 || c.DetectionMinFrames > 10 {
		return errors.New("detection_threads must be 1–16 and detection_min_frames must be 1–10")
	}
	for _, limit := range []struct {
		name            string
		value, min, max float64
	}{
		{"detection_fps", c.DetectionFPS, .5, 10},
		{"detection_confidence", c.DetectionConfidence, .05, .95},
		{"detection_merge_gap", c.DetectionMergeGap, 0, 10},
		{"detection_car_movement", c.DetectionCarMovement, .001, .25},
	} {
		if !finite(limit.value) || limit.value < limit.min || limit.value > limit.max {
			return fmt.Errorf("%s must be between %g and %g", limit.name, limit.min, limit.max)
		}
	}
	for name, value := range map[string]string{"detection_every": c.DetectionEvery, "detection_min_age": c.DetectionMinAge, "detection_timeout": c.DetectionTimeout} {
		d, err := time.ParseDuration(value)
		minimum := time.Minute
		if name == "detection_min_age" {
			minimum = 10 * time.Second
		}
		if err != nil || d < minimum || d > 24*time.Hour {
			return fmt.Errorf("%s must be between %s and 24h", name, minimum)
		}
	}
	if c.DetectionPrune != "off" && c.DetectionPrune != "quarantine" && c.DetectionPrune != "delete" {
		return errors.New("detection_prune must be off, quarantine, or delete")
	}
	return nil
}

func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func (d DetectionResult) validate() error {
	if d.Version != detectionVersion || d.Model != "yolo26n" || len(d.ModelSHA256) != 64 || len(d.Policy) != 64 || !d.Complete ||
		!finite(d.Duration) || d.Duration <= 0 || !finite(d.SampleFPS) || d.SampleFPS < .5 || d.SampleFPS > 10 || d.Frames < 1 ||
		math.Abs(float64(d.Frames)/d.SampleFPS-d.Duration) > 1.5/d.SampleFPS || d.Intervals == nil {
		return errors.New("incomplete or invalid detection analysis")
	}
	for _, interval := range d.Intervals {
		if interval.Kind != "person" && interval.Kind != "animal" && interval.Kind != "bike" && interval.Kind != "moving_car" {
			return errors.New("unknown detection category")
		}
		if !finite(interval.Start) || !finite(interval.End) || interval.Start < 0 || interval.Start >= interval.End || interval.End > d.Duration+.001 ||
			!finite(interval.Confidence) || interval.Confidence < 0 || interval.Confidence > 1 || len(interval.Labels) == 0 {
			return errors.New("invalid detection interval")
		}
	}
	return nil
}

func (s *Server) detectionPolicy() string {
	c := s.detectionConfig
	data, _ := json.Marshal([]any{"yolo26n-onnx-v1", c.DetectionModel, c.DetectionFPS, c.DetectionConfidence, c.DetectionMinFrames, c.DetectionMergeGap, c.DetectionCarMovement})
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func (s *Server) detectionCachePath(path string) string {
	h := sha256.Sum256([]byte(path))
	dir := s.detectionDir
	if dir == "" {
		dir = filepath.Join(s.root, ".neonvr_detections")
	}
	return filepath.Join(dir, hex.EncodeToString(h[:])+".json")
}

func (s *Server) cachedDetection(path string, info os.FileInfo) *DetectionResult {
	data, err := os.ReadFile(s.detectionCachePath(path))
	if err != nil {
		return nil
	}
	var cached detectionCache
	if json.Unmarshal(data, &cached) != nil || cached.Size != info.Size() || cached.MtimeNS != info.ModTime().UnixNano() ||
		cached.Result.Policy != s.detectionPolicy() || cached.Result.validate() != nil {
		return nil
	}
	return &cached.Result
}

func (s *Server) cacheDetection(path string, info os.FileInfo, result DetectionResult) error {
	data, err := json.Marshal(detectionCache{Size: info.Size(), MtimeNS: info.ModTime().UnixNano(), Result: result})
	if err != nil {
		return err
	}
	return atomicWriteConfig(s.detectionCachePath(path), data)
}

type videoInspection struct {
	Streams []struct {
		CodecType string `json:"codec_type"`
		CodecName string `json:"codec_name"`
	} `json:"streams"`
	Format struct {
		Duration string            `json:"duration"`
		Tags     map[string]string `json:"tags"`
	} `json:"format"`
}

func (s *Server) ffprobeExecutable() string {
	if strings.ContainsRune(s.ffmpegPath, filepath.Separator) {
		return filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
	}
	return "ffprobe"
}

func (s *Server) inspectDetectionVideo(ctx context.Context, path string) (videoInspection, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
		"format=duration:format_tags="+detectionMetadataTag+":stream=codec_type,codec_name", "-of", "json", path).Output()
	var info videoInspection
	if err != nil {
		return info, fmt.Errorf("ffprobe: %w", err)
	}
	if err := json.Unmarshal(out, &info); err != nil {
		return info, err
	}
	return info, nil
}

func (s *Server) embeddedDetection(ctx context.Context, path string) (*DetectionResult, error) {
	info, err := s.inspectDetectionVideo(ctx, path)
	if err != nil {
		return nil, err
	}
	data := info.Format.Tags[detectionMetadataTag]
	if data == "" {
		return nil, nil
	}
	var result DetectionResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	if err := result.validate(); err != nil {
		return nil, err
	}
	if result.Policy != s.detectionPolicy() {
		return nil, nil
	}
	return &result, nil
}

func (s *Server) apiDetectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
		return
	}
	path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	result := s.cachedDetection(path, info)
	if result == nil {
		result, err = s.embeddedDetection(r.Context(), path)
		if err == nil && result != nil {
			_ = s.cacheDetection(path, info, *result)
		}
	}
	if result != nil {
		writeJSON(w, detectionResponse{Status: "complete", Detections: result})
		return
	}
	response := detectionResponse{Status: "pending"}
	if !s.detectionConfig.DetectionEnabled {
		response.Status = "disabled"
	}
	s.detectionMu.Lock()
	if state, ok := s.detectionStates[path]; ok {
		response = state
	} else if s.detectionWorkerError != "" {
		response = detectionResponse{Status: "error", Message: s.detectionWorkerError}
	}
	s.detectionMu.Unlock()
	writeJSON(w, response)
}

type detectionWorker struct {
	cmd         *exec.Cmd
	cancel      context.CancelFunc
	stdin       io.WriteCloser
	scanner     *bufio.Scanner
	modelSHA256 string
}

func (worker *detectionWorker) close() {
	_ = worker.stdin.Close()
	worker.cancel()
	_ = worker.cmd.Wait()
}

func (worker *detectionWorker) read(ctx context.Context) ([]byte, error) {
	type reply struct {
		data []byte
		err  error
	}
	done := make(chan reply, 1)
	go func() {
		if worker.scanner.Scan() {
			done <- reply{data: append([]byte(nil), worker.scanner.Bytes()...)}
		} else {
			err := worker.scanner.Err()
			if err == nil {
				err = errors.New("detection worker exited; check server logs")
			}
			done <- reply{err: err}
		}
	}()
	select {
	case response := <-done:
		return response.data, response.err
	case <-ctx.Done():
		worker.cancel() // Kill the process group, including any FFmpeg decoder.
		<-done
		return nil, ctx.Err()
	}
}

func (s *Server) startDetectionWorker() (*detectionWorker, error) {
	c := s.detectionConfig
	ctx, cancel := context.WithCancel(s.workContext())
	cmd := exec.CommandContext(ctx, expandPath(c.DetectionPython), "-u", "-c", detectionWorkerScript,
		"--model", expandPath(c.DetectionModel), "--threads", strconv.Itoa(c.DetectionThreads))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = 2 * time.Second
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		cancel()
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		stdin.Close()
		cancel()
		return nil, err
	}
	worker := &detectionWorker{cmd: cmd, cancel: cancel, stdin: stdin, scanner: bufio.NewScanner(stdout)}
	worker.scanner.Buffer(make([]byte, 65536), 4<<20)
	readyCtx, readyCancel := context.WithTimeout(ctx, time.Minute)
	defer readyCancel()
	data, err := worker.read(readyCtx)
	var ready struct {
		Ready       bool   `json:"ready"`
		ModelSHA256 string `json:"model_sha256"`
	}
	if err == nil {
		err = json.Unmarshal(data, &ready)
	}
	if err == nil && (!ready.Ready || len(ready.ModelSHA256) != 64) {
		err = errors.New("invalid worker handshake")
	}
	if err != nil {
		worker.close()
		return nil, err
	}
	worker.modelSHA256 = ready.ModelSHA256
	return worker, nil
}

func (s *Server) analyzeDetection(ctx context.Context, worker *detectionWorker, path string) (DetectionResult, error) {
	c := s.detectionConfig
	request := map[string]any{"path": path, "ffmpeg": s.ffmpegPath, "ffprobe": s.ffprobeExecutable(),
		"fps": c.DetectionFPS, "confidence": c.DetectionConfidence, "min_frames": c.DetectionMinFrames,
		"merge_gap": c.DetectionMergeGap, "car_movement": c.DetectionCarMovement, "policy": s.detectionPolicy()}
	data, _ := json.Marshal(request)
	var result DetectionResult
	if _, err := worker.stdin.Write(append(data, '\n')); err != nil {
		return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
	}
	out, err := worker.read(ctx)
	if err != nil {
		return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
	}
	var response struct {
		Result DetectionResult `json:"result"`
		Error  string          `json:"error"`
	}
	if err := json.Unmarshal(out, &response); err != nil {
		return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
	}
	if response.Error != "" {
		return result, errors.New(response.Error)
	}
	result = response.Result
	if result.Policy != s.detectionPolicy() || result.ModelSHA256 != worker.modelSHA256 {
		return result, fmt.Errorf("%w: policy/model mismatch", errDetectionWorker)
	}
	return result, result.validate()
}

func sameRecording(a, b os.FileInfo) bool {
	return os.SameFile(a, b) && a.Size() == b.Size() && a.ModTime().Equal(b.ModTime())
}

// Remux without re-encoding. Verify metadata, streams and duration before an
// atomic replacement, retaining permissions, ownership and recording mtime.
func (s *Server) embedDetection(ctx context.Context, path string, original os.FileInfo, result DetectionResult) error {
	if err := result.validate(); err != nil {
		return err
	}
	before, err := s.inspectDetectionVideo(ctx, path)
	if err != nil {
		return err
	}
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".neonvr-detect-*.mp4")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)
	out, err := exec.CommandContext(ctx, s.ffmpegPath, "-nostdin", "-v", "error", "-y", "-i", path,
		"-map", "0", "-c", "copy", "-map_metadata", "0", "-metadata", detectionMetadataTag+"="+string(data),
		"-movflags", "+faststart+use_metadata_tags", tmpPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("detection metadata remux: %w: %.2000s", err, out)
	}
	after, err := s.inspectDetectionVideo(ctx, tmpPath)
	if err != nil {
		return err
	}
	beforeDuration, e1 := strconv.ParseFloat(before.Format.Duration, 64)
	afterDuration, e2 := strconv.ParseFloat(after.Format.Duration, 64)
	if after.Format.Tags[detectionMetadataTag] != string(data) || !reflect.DeepEqual(before.Streams, after.Streams) ||
		e1 != nil || e2 != nil || !finite(beforeDuration) || !finite(afterDuration) || math.Abs(beforeDuration-afterDuration) > .25 {
		return errors.New("remux verification failed; original recording retained")
	}
	f, err := os.OpenFile(tmpPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	if os.Geteuid() == 0 {
		if stat, ok := original.Sys().(*syscall.Stat_t); ok {
			err = f.Chown(int(stat.Uid), int(stat.Gid))
		}
	}
	if err == nil {
		err = f.Chmod(original.Mode().Perm())
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
	if err := os.Chtimes(tmpPath, original.ModTime(), original.ModTime()); err != nil {
		return err
	}
	s.clipMutationMu.Lock()
	defer s.clipMutationMu.Unlock()
	current, err := os.Stat(path)
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if !sameRecording(original, current) {
		return errors.New("recording changed during analysis; retaining it")
	}
	oldThumb, _ := s.thumbnailPath(path)
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	if oldThumb != "" {
		_ = os.Remove(oldThumb)
	}
	current, err = os.Stat(path)
	if err != nil {
		return err
	}
	s.invalidateSizeCache()
	return s.cacheDetection(path, current, result)
}

func (s *Server) setDetectionState(path, status, message string) {
	s.detectionMu.Lock()
	defer s.detectionMu.Unlock()
	if s.detectionStates == nil {
		s.detectionStates = make(map[string]detectionResponse)
	}
	if status == "complete" {
		delete(s.detectionStates, path)
		return
	}
	if len(s.detectionStates) >= 256 {
		s.detectionStates = make(map[string]detectionResponse)
	}
	s.detectionStates[path] = detectionResponse{Status: status, Message: message}
}

func (s *Server) detectionLoop() {
	if !s.waitWork(20 * time.Second) {
		return
	}
	every, _ := time.ParseDuration(s.detectionConfig.DetectionEvery)
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	var worker *detectionWorker
	defer func() {
		if worker != nil {
			worker.close()
		}
	}()
	for {
		if worker == nil {
			var err error
			worker, err = s.startDetectionWorker()
			s.detectionMu.Lock()
			s.detectionWorkerError = ""
			if err != nil {
				s.detectionWorkerError = "YOLO26n unavailable: " + err.Error()
			}
			s.detectionMu.Unlock()
			if err != nil {
				log.Printf("detection: %v; retaining recordings", err)
			}
		}
		if worker != nil {
			if err := s.runDetectionMaintenance(worker); err != nil {
				log.Printf("detection worker: %v; restarting on next pass", err)
				worker.close()
				worker = nil
			}
		}
		select {
		case <-s.workContext().Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Server) runDetectionMaintenance(worker *detectionWorker) error {
	cameras, err := s.listCameras()
	if err != nil {
		return err
	}
	minAge, _ := time.ParseDuration(s.detectionConfig.DetectionMinAge)
	timeout, _ := time.ParseDuration(s.detectionConfig.DetectionTimeout)
	for _, cam := range cameras {
		files, err := listCameraFiles(cam)
		if err != nil {
			continue
		}
		// Newest files first within each camera, so fresh events get markers early.
		sortDetectionFiles(files)
		for _, file := range files {
			if s.workContext().Err() != nil {
				return s.workContext().Err()
			}
			original, err := os.Lstat(file.FullPath)
			if err != nil || !original.Mode().IsRegular() || time.Since(original.ModTime()) < minAge {
				continue
			}
			ctx, cancel := context.WithTimeout(s.workContext(), timeout)
			result := s.cachedDetection(file.FullPath, original)
			cacheHit := result != nil
			if result == nil {
				result, _ = s.embeddedDetection(ctx, file.FullPath)
			}
			if result != nil && result.ModelSHA256 != worker.modelSHA256 {
				result = nil
			}
			if result == nil {
				s.setDetectionState(file.FullPath, "processing", "")
				analysis, analysisErr := s.analyzeDetection(ctx, worker, file.FullPath)
				if analysisErr != nil {
					s.setDetectionState(file.FullPath, "error", analysisErr.Error())
					log.Printf("detection: %s: %v; retaining recording", file.FullPath, analysisErr)
					cancel()
					// Protocol errors/timeouts invalidate the worker; a per-video
					// decoder error leaves it usable for the next recording.
					if errors.Is(analysisErr, errDetectionWorker) {
						return analysisErr
					}
					continue
				}
				if err := s.embedDetection(ctx, file.FullPath, original, analysis); err != nil {
					s.setDetectionState(file.FullPath, "error", err.Error())
					log.Printf("detection: %s: %v; retaining recording", file.FullPath, err)
					cancel()
					continue
				}
				result = &analysis
			} else if !cacheHit {
				_ = s.cacheDetection(file.FullPath, original, *result)
			}
			s.setDetectionState(file.FullPath, "complete", "")
			if err := s.pruneDetectionClip(ctx, cam, file.RelPath, *result); err != nil {
				log.Printf("detection prune: %s: %v", file.FullPath, err)
			}
			cancel()
		}
	}
	return nil
}

func sortDetectionFiles(files []clipFile) {
	sort.Slice(files, func(i, j int) bool { return files[i].Info.ModTime().After(files[j].Info.ModTime()) })
}

func (s *Server) pruneDetectionClip(ctx context.Context, cam Camera, file string, result DetectionResult) error {
	mode := s.detectionConfig.DetectionPrune
	if mode == "off" || mode == "" || result.validate() != nil || result.Policy != s.detectionPolicy() || result.Uncertain || len(result.Intervals) != 0 {
		return nil
	}
	s.clipMutationMu.Lock()
	defer s.clipMutationMu.Unlock()
	if s.meta.Get(clipKey(cam.Name, file)).Favorite {
		return nil
	}
	path, err := s.resolveClipPath(cam.Name, file)
	if err != nil {
		return err
	}
	original, err := os.Lstat(path)
	if err != nil || !original.Mode().IsRegular() {
		return errors.New("prune requires a regular recording")
	}
	// Re-read the MP4: a stale sidecar or a recording replaced during analysis
	// must never authorize removing footage.
	embedded, err := s.embeddedDetection(ctx, path)
	if err != nil {
		return err
	}
	if embedded == nil || !reflect.DeepEqual(*embedded, result) {
		return errors.New("prune metadata mismatch; retaining recording")
	}
	current, err := os.Stat(path)
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if !sameRecording(original, current) {
		return errors.New("recording changed before pruning")
	}
	thumb, _ := s.thumbnailPath(path)
	if mode == "delete" {
		err = os.Remove(path)
	} else if mode == "quarantine" {
		destination := filepath.Join(cam.Dir, ".neonvr_pruned", filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(destination), 0750); err != nil {
			return err
		}
		// Link then unlink avoids overwriting an earlier quarantined recording.
		if err := os.Link(path, destination); err != nil {
			return err
		}
		err = os.Remove(path)
		if err != nil {
			_ = os.Remove(destination)
		}
	} else {
		return errors.New("unknown prune policy")
	}
	if err != nil {
		return err
	}
	if thumb != "" {
		_ = os.Remove(thumb)
	}
	_ = os.Remove(s.detectionCachePath(path))
	s.invalidateSizeCache()
	log.Printf("detection: %s %s (fully analyzed, no targets)", mode, path)
	return nil
}
