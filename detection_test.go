package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func detectionFixture(t *testing.T) (*Server, Camera, string, DetectionResult) {
	t.Helper()
	for _, name := range []string{"ffmpeg", "ffprobe"} {
		if _, err := exec.LookPath(name); err != nil {
			t.Skip(name + " not installed")
		}
	}
	c := defaultConfig()
	c.Root = t.TempDir()
	cam := Camera{Name: "Garden", Dir: filepath.Join(c.Root, "Garden"), Managed: true}
	if err := os.MkdirAll(cam.Dir, 0750); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cam.Dir, "2026-10-01_12-00-00_Garden.mp4")
	out, err := exec.Command("ffmpeg", "-nostdin", "-v", "error", "-f", "lavfi", "-i", "color=c=green:s=160x96:r=10",
		"-f", "lavfi", "-i", "sine=frequency=440", "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1", "-c:a", "aac",
		"-metadata", "source_camera=Garden", "-metadata", "creation_time="+motionCreationTime, "-metadata", "yi_motion="+motionTestTag, "-movflags", "+use_metadata_tags", path).CombinedOutput()
	if err != nil {
		t.Fatalf("create fixture: %v: %s", err, out)
	}
	if err := os.Chmod(path, 0640); err != nil {
		t.Fatal(err)
	}
	mtime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
	meta, err := NewMetaStore(filepath.Join(c.Root, ".blinkclips_meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{root: c.Root, meta: meta, ffmpegPath: "ffmpeg", thumbDir: filepath.Join(c.Root, ".thumbs"), detectionConfig: c}
	d := DetectionResult{Version: 1, Model: "yolo26n", ModelSHA256: strings.Repeat("a", 64), Policy: s.detectionPolicy(),
		Complete: true, Duration: 4, SampleFPS: 2, Frames: 8, Intervals: []DetectionInterval{}, AnalyzedAt: time.Now().UTC().Format(time.RFC3339)}
	return s, cam, path, d
}

func streamHash(t *testing.T, path string) []byte {
	t.Helper()
	out, err := exec.Command("ffmpeg", "-nostdin", "-v", "error", "-i", path, "-map", "0", "-c", "copy", "-f", "streamhash", "-hash", "sha256", "-").CombinedOutput()
	if err != nil {
		t.Fatalf("stream hash: %v: %s", err, out)
	}
	return out
}

func TestDetectionMP4RoundTripPreservesRecording(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	result.Intervals = []DetectionInterval{
		{Kind: "person", Labels: []string{"person"}, Start: .5, End: 1.5, Confidence: .9},
		{Kind: "bike", Labels: []string{"bicycle"}, Start: 2, End: 3.5, Confidence: .8},
	}
	original, _ := os.Stat(path)
	before := streamHash(t, path)
	if err := s.embedDetection(context.Background(), path, original, result); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(path)
	if !after.ModTime().Equal(original.ModTime()) || after.Mode().Perm() != original.Mode().Perm() {
		t.Fatal("recording mtime/permissions changed")
	}
	if !bytes.Equal(before, streamHash(t, path)) {
		t.Fatal("audio/video packets changed")
	}
	metadata, err := s.embeddedDetection(context.Background(), path)
	if err != nil || !reflect.DeepEqual(metadata, &result) {
		t.Fatalf("embedded intervals did not round trip: %+v %v", metadata, err)
	}
	allTags, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format_tags=source_camera", "-of", "json", path).Output()
	if err != nil || !bytes.Contains(allTags, []byte("Garden")) {
		t.Fatal("existing custom metadata lost")
	}
	motion, err := s.readCameraMotion(context.Background(), path)
	if err != nil || motion == nil || len(motion.Samples) != 3 || motion.Samples[0].Time != 1.5 {
		t.Fatalf("camera motion metadata lost during AI remux: %+v %v", motion, err)
	}
	// The MP4 remains the source of truth when the expendable JSON cache is lost.
	if err := os.Remove(s.detectionCachePath(path)); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.apiDetectionsHandler(w, httptest.NewRequest(http.MethodGet, "/api/detections?camera=Garden&file="+filepath.Base(path), nil))
	var response detectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || response.Status != "complete" || !reflect.DeepEqual(response.Detections, &result) {
		t.Fatalf("metadata API: %d %s", w.Code, w.Body.String())
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("target recording pruned")
	}
}

func TestDetectionPruningKeepsFavoritesAndUncertainAndQuarantinesEmpty(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	result.Uncertain = true
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("uncertain recording pruned")
	}
	result.Uncertain = false
	info, _ = os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	key := clipKey(cam.Name, filepath.Base(path))
	if err := s.meta.SetFavorite(key, true); err != nil {
		t.Fatal(err)
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("favorite recording pruned")
	}
	if err := s.meta.SetFavorite(key, false); err != nil {
		t.Fatal(err)
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("empty recording not quarantined")
	}
	quarantined := filepath.Join(cam.Dir, ".neonvr_pruned", filepath.Base(path))
	metadata, err := s.embeddedDetection(context.Background(), quarantined)
	if err != nil || !reflect.DeepEqual(metadata, &result) {
		t.Fatal("quarantined video/metadata not recoverable")
	}
	files, err := listCameraFiles(cam)
	if err != nil || len(files) != 0 {
		t.Fatal("quarantine still appears in recordings")
	}
}

func TestDetectionPruneRequiresCurrentMP4Metadata(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err == nil {
		t.Fatal("pruned an MP4 with no embedded analysis")
	}
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	result.ModelSHA256 = strings.Repeat("b", 64)
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err == nil {
		t.Fatal("stale analysis authorized pruning")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("recording was removed")
	}
}

func TestDetectionRejectsChangedRecordingBeforeReplacing(t *testing.T) {
	s, _, path, result := detectionFixture(t)
	original, _ := os.Stat(path)
	newTime := original.ModTime().Add(time.Minute)
	if err := os.Chtimes(path, newTime, newTime); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if err := s.embedDetection(context.Background(), path, original, result); err == nil {
		t.Fatal("overwrote changed recording")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("failed remux changed recording")
	}
	files, _ := listCameraFiles(Camera{Dir: filepath.Dir(path)})
	if len(files) != 1 {
		t.Fatal("temporary remux leaked into recording list")
	}
}

func TestDetectionDeletionMode(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	s.detectionConfig.DetectionPrune = "delete"
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("delete mode did not delete")
	}
}

func TestDetectionCacheInvalidatesWhenPolicyOrSourceChanges(t *testing.T) {
	s, _, path, result := detectionFixture(t)
	info, _ := os.Stat(path)
	if err := s.cacheDetection(path, info, result); err != nil {
		t.Fatal(err)
	}
	if s.cachedDetection(path, info) == nil {
		t.Fatal("cache miss")
	}
	s.detectionConfig.DetectionConfidence += .1
	if s.cachedDetection(path, info) != nil {
		t.Fatal("stale policy cache used")
	}
	s.detectionConfig.DetectionConfidence -= .1
	newTime := info.ModTime().Add(time.Minute)
	os.Chtimes(path, newTime, newTime)
	changed, _ := os.Stat(path)
	if s.cachedDetection(path, changed) != nil {
		t.Fatal("changed recording cache used")
	}
}

func TestDetectionPolicyUpgradeInvalidatesOldResults(t *testing.T) {
	s, _, path, result := detectionFixture(t)
	c := s.detectionConfig
	legacy, err := json.Marshal([]any{"yolo26n-onnx-v1", c.DetectionModel, c.DetectionFPS, c.DetectionConfidence, c.DetectionMinFrames, c.DetectionMergeGap, c.DetectionCarMovement})
	if err != nil {
		t.Fatal(err)
	}
	result.Policy = configRevision(legacy)
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	info, _ = os.Stat(path)
	if s.cachedDetection(path, info) != nil {
		t.Fatal("old detection algorithm cache used")
	}
	embedded, err := s.embeddedDetection(context.Background(), path)
	if err != nil || embedded != nil {
		t.Fatalf("old embedded detection algorithm used: %+v %v", embedded, err)
	}
}

func TestDetectionConfigLimits(t *testing.T) {
	c := defaultConfig()
	for _, change := range []func(*Config){
		func(c *Config) { c.DetectionThreads = 0 }, func(c *Config) { c.DetectionFPS = 0 },
		func(c *Config) { c.DetectionConfidence = 1 }, func(c *Config) { c.DetectionModel = "yolo26n.pt" },
		func(c *Config) { c.DetectionMotionFPS = 0 }, func(c *Config) { c.DetectionMotionThreshold = 0 },
		func(c *Config) { c.DetectionMotionMinArea = 0 }, func(c *Config) { c.DetectionROIPreRoll = -1 }, func(c *Config) { c.DetectionROIPostRoll = 61 },
		func(c *Config) { c.DetectionMinFrames = 0 }, func(c *Config) { c.DetectionCarMovement = 0 },
		func(c *Config) { c.DetectionPrune = "typo" }, func(c *Config) { c.DetectionEvery = "0" },
		func(c *Config) { c.DetectionMinAge = "1s" }, func(c *Config) { c.DetectionTimeout = "0" },
	} {
		copy := c
		change(&copy)
		if copy.validateDetection() == nil {
			t.Fatal("invalid detection setting accepted")
		}
	}
}

func TestDetectionBackgroundStopsForRestart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{runCtx: ctx}
	done := make(chan struct{})
	go func() { s.detectionLoop(); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("detection did not stop")
	}
}

// Optional integration with actual weights, without making ordinary Go tests
// download models or install Python dependencies.
func TestYOLO26WorkerAndPipeline(t *testing.T) {
	model := os.Getenv("NEONVR_TEST_MODEL")
	python := os.Getenv("NEONVR_TEST_PYTHON")
	if model == "" || python == "" {
		t.Skip("set NEONVR_TEST_MODEL and NEONVR_TEST_PYTHON for real CPU inference")
	}
	s, cam, path, _ := detectionFixture(t)
	s.detectionConfig.DetectionModel = model
	s.detectionConfig.DetectionPython = python
	worker, err := s.startDetectionWorker()
	if err != nil {
		t.Fatal(err)
	}
	defer worker.close()
	bad := filepath.Join(cam.Dir, "2026-10-01_11-00-00_Corrupt.mp4")
	if err := os.WriteFile(bad, []byte("unfinished upload"), 0640); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-24 * time.Hour)
	os.Chtimes(bad, old, old)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if _, err := s.analyzeDetection(ctx, worker, bad); err == nil || errors.Is(err, errDetectionWorker) {
		t.Fatalf("expected a per-video error, got %v", err)
	}
	result, err := s.analyzeDetection(ctx, worker, path)
	if err != nil || len(result.Intervals) != 0 || result.Uncertain || result.Frames != 8 {
		t.Fatalf("CPU inference after bad video: %+v %v", result, err)
	}
	if err := s.runDetectionMaintenance(worker); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("no-target clip not pruned by real pipeline")
	}
	if _, err := os.Stat(bad); err != nil {
		t.Fatal("corrupt recording removed")
	}
	if _, err := os.Stat(filepath.Join(cam.Dir, ".neonvr_pruned", filepath.Base(path))); err != nil {
		t.Fatal("real pipeline did not retain recoverable recording")
	}
}

func TestTwoPassCoverageValidation(t *testing.T) {
	base := DetectionResult{Version: 2, Model: "yolo26m", ModelSHA256: strings.Repeat("a", 64), Policy: strings.Repeat("b", 64), Complete: true,
		Duration: 24, SampleFPS: 2, Frames: 20, MotionFPS: 5, MotionFrames: 120,
		Regions:   []DetectionRegion{{Start: 0, End: 4, Frames: 8}, {Start: 10, End: 16, Frames: 12}},
		Intervals: []DetectionInterval{{Kind: "person", Labels: []string{"person"}, Start: 11, End: 12, Confidence: .9}}}
	if err := base.validate(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		change func(*DetectionResult)
	}{
		{"missing motion coverage", func(d *DetectionResult) { d.MotionFrames = 0 }},
		{"truncated motion pass", func(d *DetectionResult) { d.MotionFrames = 20 }},
		{"missing opening window", func(d *DetectionResult) { d.Regions[0].Start = 1 }},
		{"missing YOLO region", func(d *DetectionResult) { d.Regions = nil }},
		{"truncated YOLO region", func(d *DetectionResult) { d.Regions[1].Frames = 1; d.Frames = 9 }},
		{"mismatched total", func(d *DetectionResult) { d.Frames = 21 }},
		{"overlapping regions", func(d *DetectionResult) { d.Regions[1].Start = 3 }},
		{"out of bounds", func(d *DetectionResult) { d.Regions[1].End = 25 }},
		{"detection in unexamined gap", func(d *DetectionResult) { d.Intervals[0].Start = 6; d.Intervals[0].End = 7 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := base
			d.Regions = append([]DetectionRegion(nil), base.Regions...)
			d.Intervals = append([]DetectionInterval(nil), base.Intervals...)
			tc.change(&d)
			if d.validate() == nil {
				t.Fatal("invalid two-pass analysis accepted")
			}
		})
	}
}

func TestMotionSettingsInvalidatePolicy(t *testing.T) {
	c := defaultConfig()
	s := &Server{detectionConfig: c}
	before := s.detectionPolicy()
	for _, change := range []func(*Config){
		func(c *Config) { c.DetectionMotionFPS = 10 }, func(c *Config) { c.DetectionMotionThreshold = 20 },
		func(c *Config) { c.DetectionMotionMinArea = .002 }, func(c *Config) { c.DetectionROIPreRoll = 3 }, func(c *Config) { c.DetectionROIPostRoll = 4 },
	} {
		s.detectionConfig = c
		change(&s.detectionConfig)
		if before == s.detectionPolicy() {
			t.Fatal("motion settings did not invalidate old analysis")
		}
	}
}

func TestTwoPassMP4MetadataRoundTrip(t *testing.T) {
	s, _, path, result := detectionFixture(t)
	result.Version = 2
	result.Model = "yolo26m"
	result.MotionFPS = 5
	result.MotionFrames = 20
	result.Regions = []DetectionRegion{{Start: 0, End: 4, Frames: 8}}
	result.Intervals = []DetectionInterval{{Kind: "person", Labels: []string{"person"}, Start: 1, End: 2, Confidence: .9}}
	original, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, original, result); err != nil {
		t.Fatal(err)
	}
	embedded, err := s.embeddedDetection(context.Background(), path)
	if err != nil || !reflect.DeepEqual(embedded, &result) {
		t.Fatalf("two-pass metadata lost: %+v %v", embedded, err)
	}
}
