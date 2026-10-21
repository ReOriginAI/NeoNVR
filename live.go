package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type LiveStream struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	BackchannelURL string `json:"backchannel_url,omitempty"`
}

func validateRTSP(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "rtsp" && u.Scheme != "rtsps") || u.Hostname() == "" || u.Fragment != "" || strings.ContainsAny(raw, "\r\n\t ") {
		return errors.New("enter an rtsp:// or rtsps:// camera URL without spaces or a fragment")
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return errors.New("RTSP port must be between 1 and 65535")
		}
	}
	return nil
}

func (c Config) validateLive() error {
	if strings.TrimSpace(c.LiveGo2RTC) == "" {
		return errors.New("live_go2rtc cannot be empty")
	}
	if c.LiveWebRTCPort < 1 || c.LiveWebRTCPort > 65535 {
		return errors.New("live_webrtc_port must be between 1 and 65535")
	}
	_, appPort, _ := net.SplitHostPort(c.Addr)
	if appPort == strconv.Itoa(c.LiveWebRTCPort) {
		return errors.New("live_webrtc_port must differ from the app listen port")
	}
	if host := c.LiveWebRTCHost; host != "" && (strings.ContainsAny(host, "/?#@ \r\n\t") || strings.Contains(host, ":") && net.ParseIP(host) == nil) {
		return errors.New("live_webrtc_host must be a server hostname or IP address without a port")
	}
	if len(c.LiveStreams) > 32 {
		return errors.New("configure at most 32 live streams")
	}
	names := map[string]bool{}
	for index, stream := range c.LiveStreams {
		name := strings.TrimSpace(stream.Name)
		if name == "" || len(name) > 120 || strings.ContainsAny(name, "|\r\n") || names[name] {
			return fmt.Errorf("live stream %d needs a unique name of up to 120 characters", index+1)
		}
		names[name] = true
		if err := validateRTSP(stream.URL); err != nil {
			return fmt.Errorf("live stream %d: %w", index+1, err)
		}
		if stream.BackchannelURL != "" {
			if err := validateRTSP(stream.BackchannelURL); err != nil {
				return fmt.Errorf("live stream %d backchannel: %w", index+1, err)
			}
		}
	}
	return nil
}

type liveService struct {
	cmd      *exec.Cmd
	done     chan struct{}
	baseURL  string
	password string
	client   *http.Client
}

type liveManager struct {
	mu      sync.Mutex
	ctx     context.Context
	config  Config
	token   string
	service *liveService
}

func newLiveManager(ctx context.Context, config Config, token string) *liveManager {
	// Enabling is deliberately transient: every app start defaults to off.
	return &liveManager{ctx: ctx, config: config, token: token}
}

func (m *liveManager) runningLocked() *liveService {
	if m.service != nil {
		select {
		case <-m.service.done:
			m.service = nil
		default:
		}
	}
	return m.service
}

func (m *liveManager) response(w http.ResponseWriter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	streams := make([]map[string]any, 0, len(m.config.LiveStreams))
	for index, stream := range m.config.LiveStreams {
		u, _ := url.Parse(stream.URL)
		streams = append(streams, map[string]any{"id": index, "name": stream.Name, "host": u.Host, "talk": stream.BackchannelURL != ""})
	}
	writeJSON(w, map[string]any{"enabled": m.runningLocked() != nil, "streams": streams, "token": m.token})
}

func (m *liveManager) authorize(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-NeoNVR-Token") != m.token {
		http.Error(w, "reload Streams before connecting", http.StatusForbidden)
		return false
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		// Allow the HTTPS origin used by a TLS reverse proxy on the same host.
		if err != nil || u.Host != r.Host || (u.Scheme != "http" && u.Scheme != "https") {
			http.Error(w, "same-origin requests required", http.StatusForbidden)
			return false
		}
	}
	return true
}

func (m *liveManager) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	switch r.Method {
	case http.MethodGet:
		m.response(w)
	case http.MethodPost:
		if !m.authorize(w, r) {
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "JSON required", http.StatusUnsupportedMediaType)
			return
		}
		var request struct {
			Enabled *bool `json:"enabled"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&request); err != nil || request.Enabled == nil || dec.Decode(new(any)) != io.EOF {
			http.Error(w, "expected an enabled boolean", http.StatusBadRequest)
			return
		}
		m.mu.Lock()
		var err error
		if *request.Enabled {
			if len(m.config.LiveStreams) == 0 {
				err = errors.New("add a camera in Settings before enabling live viewing")
			} else if m.runningLocked() == nil {
				m.service, err = m.start(r.Context(), r.Host)
			}
		} else {
			m.stopLocked()
		}
		m.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		m.response(w)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func liveSource(raw string, backchannel bool) string {
	u, _ := url.Parse(raw)
	flag := "0"
	if backchannel {
		flag = "1"
	}
	// The fragment is a go2rtc option; camera query parameters stay intact.
	u.Fragment = "backchannel=" + flag + "&timeout=10"
	return u.String()
}

func (m *liveManager) gatewayConfig(address, rtspAddress, password, host string) map[string]any {
	sources := map[string][]string{}
	for index, stream := range m.config.LiveStreams {
		name := fmt.Sprintf("camera_%d", index)
		// Convert unsupported camera audio on demand, using the shared local
		// stream so FFmpeg does not open an additional camera connection.
		sources[name] = []string{liveSource(stream.URL, false), "ffmpeg:" + name + "#audio=opus"}
		if stream.BackchannelURL != "" {
			talkName := name + "_talk"
			sources[talkName] = []string{liveSource(stream.BackchannelURL, true), "ffmpeg:" + talkName + "#audio=opus"}
		}
	}
	return map[string]any{
		"api":     map[string]any{"listen": address, "username": "neonvr", "password": password, "local_auth": true, "allow_paths": []string{"/api", "/api/webrtc", "/api/streams"}},
		"rtsp":    map[string]any{"listen": rtspAddress},
		"ffmpeg":  map[string]any{"bin": expandPath(m.config.FFmpeg)},
		"webrtc":  map[string]any{"listen": fmt.Sprintf(":%d", m.config.LiveWebRTCPort), "candidates": []string{net.JoinHostPort(host, strconv.Itoa(m.config.LiveWebRTCPort))}, "ice_servers": []any{}, "filters": map[string]any{"loopback": true}},
		"streams": sources,
		"log":     map[string]any{"level": "error"},
	}
}

func (m *liveManager) start(requestCtx context.Context, requestHost string) (*liveService, error) {
	mediaAddress := fmt.Sprintf(":%d", m.config.LiveWebRTCPort)
	tcp, err := net.Listen("tcp", mediaAddress)
	if err != nil {
		return nil, errors.New("WebRTC TCP port is already in use; choose another port in Settings")
	}
	udp, err := net.ListenPacket("udp", mediaAddress)
	tcp.Close()
	if err != nil {
		return nil, errors.New("WebRTC UDP port is already in use; choose another port in Settings")
	}
	udp.Close()
	// Reserve a private API address. Only fixed, internal API routes are used.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, errors.New("cannot allocate a local streaming gateway port")
	}
	address := listener.Addr().String()
	rtspListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		listener.Close()
		return nil, errors.New("cannot allocate a local audio streaming port")
	}
	rtspAddress := rtspListener.Addr().String()
	rtspListener.Close()
	listener.Close()
	host := m.config.LiveWebRTCHost
	if host == "" {
		host, _, err = net.SplitHostPort(requestHost)
		if err != nil {
			host = requestHost
		}
	}
	directory, err := os.MkdirTemp("", "neonvr-live-*")
	if err != nil {
		return nil, errors.New("cannot create temporary streaming configuration")
	}
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(directory)
		}
	}()
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	password := hex.EncodeToString(secret)
	data, _ := json.Marshal(m.gatewayConfig(address, rtspAddress, password, host))
	path := filepath.Join(directory, "go2rtc.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, errors.New("cannot write temporary streaming configuration")
	}
	cmd := exec.CommandContext(m.ctx, expandPath(m.config.LiveGo2RTC), "-config", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
	// RTSP errors can contain camera credentials; do not expose gateway logs.
	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	if err := cmd.Start(); err != nil {
		return nil, errors.New("cannot start go2rtc; install it or set its executable path in Settings")
	}
	service := &liveService{cmd: cmd, done: make(chan struct{}), baseURL: "http://" + address, password: password, client: &http.Client{Timeout: 25 * time.Second, Transport: &http.Transport{}}}
	cleanup = false
	go func() {
		_ = cmd.Wait()
		service.client.CloseIdleConnections()
		os.RemoveAll(directory)
		close(service.done)
	}()
	ctx, cancel := context.WithTimeout(requestCtx, 5*time.Second)
	defer cancel()
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, service.baseURL+"/api", nil)
		req.SetBasicAuth("neonvr", password)
		response, requestErr := service.client.Do(req)
		if requestErr == nil {
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return service, nil
			}
		}
		select {
		case <-service.done:
			return nil, errors.New("streaming gateway exited; check the executable and WebRTC port in Settings")
		case <-ctx.Done():
			service.stop()
			return nil, errors.New("streaming gateway did not become ready")
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (s *liveService) stop() {
	_ = s.cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-s.done:
	case <-time.After(2 * time.Second):
		_ = s.cmd.Process.Kill()
		<-s.done
	}
}

func (m *liveManager) stopLocked() {
	if service := m.service; service != nil {
		m.service = nil
		service.stop()
	}
}

func (m *liveManager) close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopLocked()
}

func (s *liveService) backchannelAccepted(ctx context.Context, source, sessionID string) bool {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/streams?src="+source, nil)
	req.SetBasicAuth("neonvr", s.password)
	response, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	var status struct {
		Consumers []struct {
			UserAgent string            `json:"user_agent"`
			Receivers []json.RawMessage `json:"receivers"`
		} `json:"consumers"`
	}
	if response.StatusCode != http.StatusOK || json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(&status) != nil {
		return false
	}
	for _, consumer := range status.Consumers {
		// A receiver exists only when the gateway attached the browser's
		// microphone to a compatible camera speaker track. SDP direction alone
		// can remain sendonly even when the camera did not accept a backchannel.
		if consumer.UserAgent == sessionID && len(consumer.Receivers) > 0 {
			return true
		}
	}
	return false
}

func (m *liveManager) webrtcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !m.authorize(w, r) {
		return
	}
	if r.Header.Get("Content-Type") != "application/sdp" {
		http.Error(w, "SDP offer required", http.StatusUnsupportedMediaType)
		return
	}
	m.mu.Lock()
	service := m.runningLocked()
	m.mu.Unlock()
	if service == nil {
		http.Error(w, "live viewing is off", http.StatusConflict)
		return
	}
	index, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || index < 0 || index >= len(m.config.LiveStreams) {
		http.Error(w, "unknown camera", http.StatusNotFound)
		return
	}
	source := fmt.Sprintf("camera_%d", index)
	talk := r.URL.Query().Get("talk") == "1"
	if talk {
		if m.config.LiveStreams[index].BackchannelURL == "" {
			http.Error(w, "this camera has no backchannel URL", http.StatusBadRequest)
			return
		}
		source += "_talk"
	}
	// Never forward arbitrary URLs or other go2rtc query parameters.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil || !bytes.HasPrefix(body, []byte("v=0")) {
		http.Error(w, "invalid SDP offer", http.StatusBadRequest)
		return
	}
	req, _ := http.NewRequestWithContext(r.Context(), http.MethodPost, service.baseURL+"/api/webrtc?src="+source, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/sdp")
	req.SetBasicAuth("neonvr", service.password)
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		http.Error(w, "cannot create a live session", http.StatusInternalServerError)
		return
	}
	sessionID := "NeoNVR/" + hex.EncodeToString(nonce)
	req.Header.Set("User-Agent", sessionID)
	response, err := service.client.Do(req)
	if err != nil {
		http.Error(w, "camera connection failed; check the camera address and network", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		// The gateway's raw errors may include RTSP credentials.
		http.Error(w, "camera could not provide compatible video or audio; check its URL, availability, and codecs", http.StatusBadGateway)
		return
	}
	if talk && !service.backchannelAccepted(r.Context(), source, sessionID) {
		http.Error(w, "camera did not accept a compatible audio backchannel; your microphone is off, choose Watch to return to viewing", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/sdp")
	_, _ = io.Copy(w, io.LimitReader(response.Body, 128<<10))
}
