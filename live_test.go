package main

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func liveFixture() Config {
	c := defaultConfig()
	c.LiveStreams = []LiveStream{{Name: "Front camera", URL: "rtsp://user:camera-secret@192.168.1.150:554/ch0_0.h264", BackchannelURL: "rtsp://user:camera-secret@192.168.1.150:554/ch0_0.h264?backchannel=1"}}
	return c
}

func liveRequest(method, path, body, contentType string) *http.Request {
	r := httptest.NewRequest(method, "http://neonvr.test"+path, strings.NewReader(body))
	r.Header.Set("X-NeoNVR-Token", "test-token")
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Origin", "http://neonvr.test")
	return r
}

func TestLiveDefaultsOffAndHidesCameraCredentials(t *testing.T) {
	c := liveFixture()
	c.LiveGo2RTC = "/does-not-exist"
	m := newLiveManager(context.Background(), c, "test-token")
	defer m.close()
	w := httptest.NewRecorder()
	m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
	var status struct {
		Enabled bool `json:"enabled"`
		Streams []struct {
			Name string `json:"name"`
			Talk bool   `json:"talk"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil || status.Enabled || len(status.Streams) != 1 || !status.Streams[0].Talk {
		t.Fatalf("unexpected default status: %s", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "camera-secret") || strings.Contains(w.Body.String(), "rtsp://") || m.service != nil {
		t.Fatal("listing sources exposed credentials or started a camera gateway")
	}
	w = httptest.NewRecorder()
	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0", "v=0\r\n", "application/sdp"))
	if w.Code != http.StatusConflict || m.service != nil {
		t.Fatalf("a direct playback request bypassed the off switch: %d", w.Code)
	}
	w = httptest.NewRecorder()
	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
	if w.Code != http.StatusServiceUnavailable || m.service != nil {
		t.Fatal("a failed gateway start left live viewing enabled")
	}
}

func TestLiveToggleRejectsUnauthorizedAndMalformedRequests(t *testing.T) {
	m := newLiveManager(context.Background(), liveFixture(), "test-token")
	defer m.close()
	for _, tc := range []struct {
		name, method, body, token, origin, contentType string
		code                                           int
	}{
		{"missing token", "POST", `{"enabled":false}`, "", "http://neonvr.test", "application/json", 403},
		{"foreign origin", "POST", `{"enabled":false}`, "test-token", "https://other.test", "application/json", 403},
		{"missing boolean", "POST", `{}`, "test-token", "http://neonvr.test", "application/json", 400},
		{"wrong type", "POST", `{"enabled":"true"}`, "test-token", "http://neonvr.test", "application/json", 400},
		{"trailing data", "POST", `{"enabled":false}{}`, "test-token", "http://neonvr.test", "application/json", 400},
		{"unexpected field", "POST", `{"enabled":false,"src":"rtsp://other"}`, "test-token", "http://neonvr.test", "application/json", 400},
		{"wrong content type", "POST", `{"enabled":false}`, "test-token", "http://neonvr.test", "text/plain", 415},
		{"wrong method", "PUT", "", "test-token", "http://neonvr.test", "application/json", 405},
		{"HTTPS proxy", "POST", `{"enabled":false}`, "test-token", "https://neonvr.test", "application/json", 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := liveRequest(tc.method, "/api/live", tc.body, tc.contentType)
			r.Header.Set("X-NeoNVR-Token", tc.token)
			r.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			m.handler(w, r)
			if w.Code != tc.code || m.service != nil {
				t.Fatalf("status=%d, want %d; %s", w.Code, tc.code, w.Body.String())
			}
		})
	}
}

func TestLiveSignalingUsesOnlyConfiguredSourcesAndRedactsErrors(t *testing.T) {
	var observed string
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || user != "neonvr" || password != "private-password" || r.URL.Path != "/api/webrtc" {
			t.Error("gateway authentication or route missing")
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "v=0\r\n" || r.Header.Get("Content-Type") != "application/sdp" {
			t.Error("SDP offer was changed")
		}
		observed = r.URL.RawQuery
		if r.URL.Query().Get("src") == "camera_0_talk" {
			http.Error(w, "rtsp://user:camera-secret@camera failed", 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("v=0\r\nanswer"))
	}))
	defer backend.Close()
	c := liveFixture()
	c.LiveStreams = append(c.LiveStreams, LiveStream{Name: "View only", URL: "rtsp://192.168.1.151/live"})
	m := newLiveManager(context.Background(), c, "test-token")
	m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private-password", client: backend.Client()}
	w := httptest.NewRecorder()
	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&src=rtsp://other&dst=other", "v=0\r\n", "application/sdp"))
	if w.Code != 200 || w.Body.String() != "v=0\r\nanswer" || observed != "src=camera_0" {
		t.Fatalf("unexpected signaling response: %d %s, query %s", w.Code, w.Body.String(), observed)
	}
	w = httptest.NewRecorder()
	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&talk=1", "v=0\r\n", "application/sdp"))
	if w.Code != 502 || observed != "src=camera_0_talk" || strings.Contains(w.Body.String(), "camera-secret") {
		t.Fatal("backchannel selection failed or camera credentials leaked")
	}
	for _, tc := range []struct {
		path, body string
		code       int
	}{
		{"/api/live/webrtc?id=-1", "v=0\r\n", 404},
		{"/api/live/webrtc?id=10", "v=0\r\n", 404},
		{"/api/live/webrtc?id=1&talk=1", "v=0\r\n", 400},
		{"/api/live/webrtc?id=0", "invalid", 400},
		{"/api/live/webrtc?id=0", "v=0" + strings.Repeat("x", 64<<10), 400},
	} {
		observed = ""
		w = httptest.NewRecorder()
		m.webrtcHandler(w, liveRequest("POST", tc.path, tc.body, "application/sdp"))
		if w.Code != tc.code || observed != "" {
			t.Fatalf("invalid request reached the gateway: %s %d", tc.path, w.Code)
		}
	}
}

func TestLiveConfigValidationAndLegacyCompatibility(t *testing.T) {
	c, err := decodeConfig([]byte(`{"root":"/tmp"}`))
	if err != nil || c.LiveStreams == nil || len(c.LiveStreams) != 0 || c.LiveGo2RTC != "go2rtc" || c.LiveWebRTCPort != 8555 {
		t.Fatalf("legacy config did not receive optional live defaults: %+v %v", c, err)
	}
	for _, raw := range []string{"exec:sh -c anything", "http://camera/live", "rtsp:///missing-host", "rtsp://camera:65536/live", "rtsp://camera/live#backchannel=1", "rtsp://camera/live\n"} {
		if validateRTSP(raw) == nil {
			t.Fatalf("invalid source accepted: %q", raw)
		}
	}
	c = liveFixture()
	if err := c.validateLive(); err != nil {
		t.Fatal(err)
	}
	m := newLiveManager(context.Background(), c, "test-token")
	gateway := m.gatewayConfig("127.0.0.1:12345", "127.0.0.1:12346", "secret", "192.168.1.10")
	sources := gateway["streams"].(map[string][]string)
	if !strings.HasSuffix(sources["camera_0"][0], "#backchannel=0&timeout=10") || !strings.Contains(sources["camera_0_talk"][0], "?backchannel=1#backchannel=1&timeout=10") {
		t.Fatal("viewing opened a backchannel or talk lost the camera query")
	}
	if sources["camera_0"][1] != "ffmpeg:camera_0#audio=opus" || sources["camera_0_talk"][1] != "ffmpeg:camera_0_talk#audio=opus" {
		t.Fatal("audio conversion does not reuse the shared viewing/talk stream")
	}
	if gateway["rtsp"].(map[string]any)["listen"] != "127.0.0.1:12346" {
		t.Fatal("audio restream port is not private")
	}
	c.LiveStreams[0].BackchannelURL = "exec:invalid"
	if c.validateLive() == nil {
		t.Fatal("invalid backchannel accepted")
	}
}

func TestLiveTalkRequiresAnAttachedCameraSpeakerTrack(t *testing.T) {
	var sessionID string
	var attached, wrongSession bool
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/webrtc":
			sessionID = r.UserAgent()
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("v=0\r\nanswer"))
		case "/api/streams":
			if r.URL.RawQuery != "src=camera_0_talk" {
				t.Error("backchannel verification queried an arbitrary source")
			}
			agent := sessionID
			if wrongSession {
				agent = "another-viewer"
			}
			receivers := []string{}
			if attached {
				receivers = append(receivers, "microphone")
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"consumers": []any{map[string]any{"user_agent": agent, "receivers": receivers}}})
		default:
			t.Error("unexpected private gateway route")
		}
	}))
	defer backend.Close()
	m := newLiveManager(context.Background(), liveFixture(), "test-token")
	m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private", client: backend.Client()}
	for _, tc := range []struct {
		attached, wrongSession bool
		code                   int
	}{
		{false, false, 502}, {true, true, 502}, {true, false, 200},
	} {
		attached, wrongSession = tc.attached, tc.wrongSession
		w := httptest.NewRecorder()
		m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&talk=1", "v=0\r\n", "application/sdp"))
		if w.Code != tc.code || !strings.HasPrefix(sessionID, "NeoNVR/") {
			t.Fatalf("talk accepted without its own attached speaker track: %d %s", w.Code, w.Body.String())
		}
	}
}

func TestLiveGatewayLifecycleWithoutOpeningCamera(t *testing.T) {
	binary := os.Getenv("NEONVR_TEST_GO2RTC")
	if binary == "" {
		t.Skip("set NEONVR_TEST_GO2RTC to test the real streaming gateway")
	}
	camera, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer camera.Close()
	media, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := media.Addr().(*net.TCPAddr).Port
	media.Close()
	c := defaultConfig()
	c.LiveGo2RTC = binary
	c.LiveWebRTCPort = port
	c.LiveStreams = []LiveStream{{Name: "Test camera", URL: "rtsp://" + camera.Addr().String() + "/live"}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := newLiveManager(ctx, c, "test-token")
	defer m.close()
	w := httptest.NewRecorder()
	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
	if w.Code != 200 || m.service == nil {
		t.Fatalf("gateway did not start: %d %s", w.Code, w.Body.String())
	}
	service := m.service
	configPath := service.cmd.Args[2]
	if info, err := os.Stat(configPath); err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("temporary camera configuration was not private")
	}
	w = httptest.NewRecorder()
	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
	if w.Code != 200 || m.service != service {
		t.Fatal("enabling twice started duplicate gateways")
	}
	_ = camera.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))
	if conn, err := camera.Accept(); err == nil {
		conn.Close()
		t.Fatal("enabling the gateway opened a camera before Watch")
	}
	w = httptest.NewRecorder()
	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":false}`, "application/json"))
	if w.Code != 200 || m.service != nil {
		t.Fatal("disabling did not stop the gateway")
	}
	select {
	case <-service.done:
	default:
		t.Fatal("gateway process survived disabling")
	}
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatal("temporary camera credentials survived disabling")
	}
	w = httptest.NewRecorder()
	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	cancel()
	select {
	case <-m.service.done:
	case <-time.After(3 * time.Second):
		t.Fatal("gateway survived app shutdown")
	}
	w = httptest.NewRecorder()
	m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
	var status struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil || status.Enabled {
		t.Fatal("exited gateway still reported as enabled")
	}
}
