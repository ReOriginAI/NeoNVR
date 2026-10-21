package main

import (
	"context"
	"encoding/json"
	"errors"
	"image/jpeg"
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

const motionCreationTime = "2026-10-12T22:27:00Z"
const motionTestTag = `{"schema":["unix_ms","Scene","Move","MovingLevel","BinImgRatio","MovingTh"],"data":[[1791844021500,0,1,1,2.10,20],[1791844021700,0,2,2,4.25,20],[1791844022500,0,2,2,3.80,20]]}`

func TestCompactCameraSamples(t *testing.T) {
	result, err := parseCameraMotion(motionTestTag, 4, motionCreationTime)
	if err != nil {
		t.Fatal(err)
	}
	want := []motionSample{{1.5, 1791844021500, 0, 1, 1, 2.1, 20}, {1.7, 1791844021700, 0, 2, 2, 4.25, 20}, {2.5, 1791844022500, 0, 2, 2, 3.8, 20}}
	if !reflect.DeepEqual(result.Samples, want) || len(result.Data) != 3 || result.InvalidRows != 0 || result.OutsideClipRows != 0 || result.CreationTime != motionCreationTime {
		t.Fatalf("wrong compact samples: %+v", result)
	}
	// A timezone offset in creation_time expresses the same UTC origin. The
	// filename's local time and the first retained row are never origins.
	shifted, err := parseCameraMotion(motionTestTag, 4, "2026-10-12T18:27:00-04:00")
	if err != nil || !reflect.DeepEqual(shifted.Samples, want) {
		t.Fatalf("UTC offsets shifted samples: %+v %v", shifted, err)
	}
	precise, err := parseCameraMotion(motionTestTag, 4, "2026-10-12T22:27:00.125500Z")
	if err != nil || precise.Samples[0].Time != 1.3745 {
		t.Fatalf("creation_time fractional precision: %+v %v", precise, err)
	}
}

func TestCompactSampleColumnOrderBoundsAndRawPreservation(t *testing.T) {
	tag := `{"schema":["MovingTh","BinImgRatio","MovingLevel","Move","Scene","unix_ms","future"],"data":[[20,3.8,2,2,0,1791844022500,"retain"],[10,0,0,0,0,1791844021500,{"raw":true}],[31,100,3,4,2,1791844021700,null],[20,1,1,1,0,1791844021500,"duplicate"],[20,null,1,1,0,1791844023000,null],[20,101,1,1,0,1791844023000,null],[9,1,1,1,0,1791844023000,null],[20,1,4,1,0,1791844023000,null],[20,1,1,5,0,1791844023000,null],[20,1,1,1,3,1791844023000,null],[20,1,1.5,1,0,1791844023000,null],[20,"1",1,1,0,1791844023000,null],null,[20,1,1,1,0,1791844019900,null],[20,1,1,1,0,1791844024000,null]]}`
	result, err := parseCameraMotion(tag, 4, motionCreationTime)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Samples) != 4 || result.Samples[0].Time != 1.5 || result.Samples[1].Time != 1.5 || result.Samples[2].MovingLevel != 3 || result.Samples[2].BinImgRatio != 100 || result.InvalidRows != 9 || result.OutsideClipRows != 2 {
		t.Fatalf("bounds or sorted sample times: %+v", result)
	}
	if len(result.Data) != 15 || string(result.Data[0][6]) != `"retain"` || string(result.Data[1][6]) != `{"raw":true}` || result.Schema[6] != "future" {
		t.Fatal("raw vendor data was changed")
	}
}

func TestCompactEmptyDataAndRejectedMetadata(t *testing.T) {
	empty := `{"schema":["unix_ms","Scene","Move","MovingLevel","BinImgRatio","MovingTh"],"data":[]}`
	result, err := parseCameraMotion(empty, 4, "")
	if err != nil || result == nil || len(result.Samples) != 0 || result.Data == nil || result.Samples == nil {
		t.Fatalf("empty data should not infer motion: %+v %v", result, err)
	}
	if _, err := parseCameraMotion(`{"schema":"yi-motion-v1","events":[]}`, 4, motionCreationTime); !errors.Is(err, errUnsupportedCameraMotion) {
		t.Fatal("old event parser is still active", err)
	}
	for _, tag := range []string{"broken", "null", `{"schema":["unix_ms"],"data":[]}`, `{"schema":["unix_ms","unix_ms"],"data":[]}`, strings.Replace(empty, `"data":[]`, `"data":null`, 1), strings.Replace(empty, `"data":[]`, `"data":{}`, 1)} {
		if _, err := parseCameraMotion(tag, 4, motionCreationTime); err == nil {
			t.Fatalf("accepted invalid compact data: %s", tag)
		}
	}
	if _, err := parseCameraMotion(motionTestTag, 0, motionCreationTime); err == nil {
		t.Fatal("accepted missing duration")
	}
	for _, origin := range []string{"", "2026-10-12T22:27:00", "broken"} {
		if _, err := parseCameraMotion(motionTestTag, 4, origin); err == nil {
			t.Fatalf("fabricated a time origin from %q", origin)
		}
	}
}

func cameraMotionFixture(t *testing.T, tag string) (*Server, string) {
	t.Helper()
	return cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=40:y=20:w=80:h=60:color=red:t=fill:enable='between(t,1.5,2.5)'")
}
func cameraMotionFixtureWithFilter(t *testing.T, tag, filter string) (*Server, string) {
	t.Helper()
	for _, tool := range []string{"ffmpeg", "ffprobe"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skip(tool + " not installed")
		}
	}
	root := t.TempDir()
	dir := filepath.Join(root, "Garden")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "recording.mp4")
	args := []string{"-nostdin", "-v", "error", "-f", "lavfi", "-i", filter, "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1", "-metadata", "creation_time=" + motionCreationTime}
	if tag != "" {
		args = append(args, "-metadata", "yi_motion="+tag)
	}
	args = append(args, "-movflags", "+use_metadata_tags", path)
	if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
		t.Fatalf("create fixture: %v: %s", err, out)
	}
	return &Server{root: root, ffmpegPath: "ffmpeg", thumbDir: filepath.Join(root, ".thumbs"), thumbWidth: 160, thumbSeek: .2, thumbSem: make(chan struct{}, 1)}, path
}

func TestCompactCameraAPIThumbnailAndMetadataPreservation(t *testing.T) {
	s, path := cameraMotionFixture(t, motionTestTag)
	w := httptest.NewRecorder()
	s.apiCameraMotionHandler(w, httptest.NewRequest("GET", "/api/camera-motion?camera=Garden&file=recording.mp4", nil))
	var response struct {
		Motion *cameraMotion `json:"motion"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != http.StatusOK || response.Motion == nil || len(response.Motion.Samples) != 3 || response.Motion.Samples[0].Time != 1.5 {
		t.Fatalf("compact MP4 sample extraction: %d %s", w.Code, w.Body.String())
	}
	thumb, err := s.ensureThumbnail(path)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(thumb)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	frame, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := frame.At(80, 48).RGBA()
	if r < 45000 || g > 12000 || b > 12000 {
		t.Fatalf("thumbnail did not select the sampled frame: %d %d %d", r, g, b)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	result := DetectionResult{Version: 1, Model: "yolo26n", ModelSHA256: strings.Repeat("a", 64), Policy: s.detectionPolicy(), Complete: true, Duration: 4, SampleFPS: 2, Frames: 8, Intervals: []DetectionInterval{}, AnalyzedAt: time.Now().UTC().Format(time.RFC3339)}
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	after, err := s.readCameraMotion(context.Background(), path)
	if err != nil || !reflect.DeepEqual(after, response.Motion) {
		t.Fatalf("AI remux changed camera readings or UTC origin: %+v %v", after, err)
	}
	for _, request := range []*http.Request{httptest.NewRequest("POST", "/api/camera-motion?camera=Garden&file=recording.mp4", nil), httptest.NewRequest("GET", "/api/camera-motion?camera=Garden&file=../outside.mp4", nil)} {
		w := httptest.NewRecorder()
		s.apiCameraMotionHandler(w, request)
		if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusNotFound {
			t.Fatalf("invalid request accepted: %d", w.Code)
		}
	}
}

func TestCameraCreationTimeVideoStreamFallback(t *testing.T) {
	// The format tag is missing, but the actual video stream supplies the origin.
	probe := map[string]any{"format": map[string]any{"duration": "4", "tags": map[string]string{"yi_motion": motionTestTag}}, "streams": []any{map[string]any{"codec_type": "audio", "tags": map[string]string{"creation_time": "2020-01-01T00:00:00Z"}}, map[string]any{"codec_type": "video", "tags": map[string]string{"creation_time": motionCreationTime}}}}
	data, _ := json.Marshal(probe)
	dir := t.TempDir()
	script := "#!/bin/sh\ncat <<'JSON'\n" + string(data) + "\nJSON\n"
	if err := os.WriteFile(filepath.Join(dir, "ffprobe"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	s := &Server{ffmpegPath: filepath.Join(dir, "ffmpeg")}
	motion, err := s.readCameraMotion(context.Background(), "recording.mp4")
	if err != nil || motion == nil || motion.Samples[0].Time != 1.5 {
		t.Fatalf("video creation_time fallback: %+v %v", motion, err)
	}
}

func TestCameraSampleThumbnailFallbackAndOldTag(t *testing.T) {
	empty := `{"schema":["unix_ms","Scene","Move","MovingLevel","BinImgRatio","MovingTh"],"data":[]}`
	for _, tag := range []string{"", "invalid-json", empty, `{"schema":"yi-motion-v1","events":[]}`} {
		s, path := cameraMotionFixture(t, tag)
		thumb, err := s.ensureThumbnail(path)
		if err != nil {
			t.Fatal(err)
		}
		if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
			t.Fatalf("fallback produced no usable image: %v", err)
		}
		if tag == "" {
			if motion, err := s.readCameraMotion(context.Background(), path); err != nil || motion != nil {
				t.Fatalf("missing tag: %+v %v", motion, err)
			}
		}
		if strings.Contains(tag, "yi-motion-v1") {
			w := httptest.NewRecorder()
			s.apiCameraMotionHandler(w, httptest.NewRequest("GET", "/api/camera-motion?camera=Garden&file=recording.mp4", nil))
			var response struct {
				Unsupported bool          `json:"unsupported"`
				Motion      *cameraMotion `json:"motion"`
			}
			if json.Unmarshal(w.Body.Bytes(), &response) != nil || w.Code != http.StatusOK || !response.Unsupported || response.Motion != nil {
				t.Fatalf("legacy tag was parsed or broke playback: %d %s", w.Code, w.Body.String())
			}
		}
	}
}

func TestCameraSampleBlankFrameFallback(t *testing.T) {
	for _, tag := range []string{
		`{"schema":["unix_ms","Scene","Move","MovingLevel","BinImgRatio","MovingTh"],"data":[[1791844020500,0,1,1,2,20]]}`,
		`{"schema":["unix_ms","Scene","Move","MovingLevel","BinImgRatio","MovingTh"],"data":[[1791844020500,0,1,1,2,20],[1791844022000,0,1,1,2,20]]}`,
	} {
		s, path := cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=0:y=0:w=iw:h=ih:color=black:t=fill:enable='lt(t,1)'")
		thumb, err := s.ensureThumbnail(path)
		if err != nil {
			t.Fatal(err)
		}
		if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
			t.Fatalf("blank camera reading frame won over later sample or scene fallback: %v", err)
		}
	}
}
