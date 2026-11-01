package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetectionRangesSurviveFreshServerAndCacheLoss(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	result.Version = 3
	result.Model = "yolo26m"
	result.MotionFPS = result.SampleFPS
	result.MotionFrames = result.Frames
	result.Regions = []DetectionRegion{{Start: 0, End: 4, Frames: 8}}
	result.Intervals = []DetectionInterval{{Kind: "person", Labels: []string{"person"}, Start: 1, End: 2, Confidence: .9}}
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	fresh := func(config Config) *Server {
		meta, err := NewMetaStore(filepath.Join(s.root, ".blinkclips_meta.json"))
		if err != nil {
			t.Fatal(err)
		}
		return &Server{root: s.root, meta: meta, ffmpegPath: s.ffmpegPath, detectionConfig: config, detectionDir: s.detectionDir}
	}
	read := func(server *Server) detectionResponse {
		t.Helper()
		w := httptest.NewRecorder()
		server.apiDetectionsHandler(w, httptest.NewRequest(http.MethodGet, "/api/detections?camera="+cam.Name+"&file="+filepath.Base(path), nil))
		var response detectionResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 {
			t.Fatalf("API: %d %s", w.Code, w.Body.String())
		}
		if response.Status != "complete" || !reflect.DeepEqual(response.Detections, &result) {
			t.Fatalf("saved ROI or YOLO ranges lost: %s", w.Body.String())
		}
		return response
	}
	if read(fresh(s.detectionConfig)).Stale {
		t.Fatal("restart invalidated analysis")
	}
	if err := os.Remove(s.detectionCachePath(path)); err != nil {
		t.Fatal(err)
	}
	if read(fresh(s.detectionConfig)).Stale {
		t.Fatal("embedded ranges not recovered after cache loss")
	}
	if _, err := os.Stat(s.detectionCachePath(path)); err != nil {
		t.Fatal("persistent cache not rebuilt from MP4")
	}
	changed := s.detectionConfig
	changed.DetectionPersonConfidence = .8
	server := fresh(changed)
	response := read(server)
	if !response.Stale || response.Message == "" {
		t.Fatal("historical analysis not labeled")
	}
	// A saved result may be shown, but must not authorize pruning under new settings.
	info, _ = os.Stat(path)
	if server.cachedDetection(path, info) != nil {
		t.Fatal("historical cache authorized current-policy analysis")
	}
	embedded, err := server.embeddedDetection(context.Background(), path)
	if err != nil || embedded != nil {
		t.Fatal("historical MP4 authorized current-policy analysis")
	}
	if err := os.Remove(s.detectionCachePath(path)); err != nil {
		t.Fatal(err)
	}
	if !read(fresh(changed)).Stale {
		t.Fatal("historical MP4 not recovered after cache loss")
	}
	// When a fresh analysis is embedded, it replaces the saved view automatically.
	result.Policy = server.detectionPolicy()
	info, _ = os.Stat(path)
	if err := server.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	if read(fresh(changed)).Stale {
		t.Fatal("refreshed analysis still historical")
	}
}
