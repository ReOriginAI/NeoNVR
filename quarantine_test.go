package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestQuarantineBrowsingPlaybackAndFavoriteProtection(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.apiQuarantineHandler(w, httptest.NewRequest(http.MethodGet, "/api/quarantine", nil))
	var response clipsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || len(response.Clips) != 1 {
		t.Fatalf("quarantine API: %d %s", w.Code, w.Body.String())
	}
	clip := response.Clips[0]
	if clip.File != ".neonvr_pruned/"+filepath.Base(path) {
		t.Fatalf("wrong path: %s", clip.File)
	}
	w = httptest.NewRecorder()
	s.videoHandler(w, httptest.NewRequest(http.MethodGet, clip.VideoURL, nil))
	if w.Code != 200 || w.Body.Len() == 0 {
		t.Fatalf("quarantine playback: %d", w.Code)
	}
	w = httptest.NewRecorder()
	s.apiDetectionsHandler(w, httptest.NewRequest(http.MethodGet, "/api/detections?camera="+url.QueryEscape(cam.Name)+"&file="+url.QueryEscape(clip.File), nil))
	var detections detectionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &detections); err != nil || detections.Status != "complete" {
		t.Fatalf("quarantine detection metadata: %s", w.Body.String())
	}
	regular, err := s.scanClips(nil, "", false)
	if err != nil || len(regular) != 0 {
		t.Fatalf("quarantine leaked into main library: %+v %v", regular, err)
	}
	if err := s.meta.SetFavorite(clip.Key, true); err != nil {
		t.Fatal(err)
	}
	favorites, err := s.scanClips(nil, "", true)
	if err != nil || len(favorites) != 1 {
		t.Fatalf("quarantine favorite not listed: %+v %v", favorites, err)
	}
	s.runStorageCleanup(func(string) (float64, error) { return 1, nil })
	if _, err := os.Stat(filepath.Join(cam.Dir, filepath.FromSlash(clip.File))); err != nil {
		t.Fatal("quarantine favorite deleted")
	}
	filtered, err := s.scanClipsInView(map[string]bool{"Other": true}, "", false, true)
	if err != nil || len(filtered) != 0 {
		t.Fatal("quarantine camera filter ignored")
	}
}

func TestQuarantineListingYiAndSymlinks(t *testing.T) {
	s, cam, _, _ := detectionFixture(t)
	dir := filepath.Join(cam.Dir, ".neonvr_pruned", "2020Y04M17D16H")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "16M36S06.mp4"), []byte("clip"), 0640); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "outside.mp4"), []byte("clip"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "outside.mp4"), filepath.Join(dir, "link.mp4")); err != nil {
		t.Fatal(err)
	}
	clips, err := s.scanClipsInView(nil, "", false, true)
	if err != nil || len(clips) != 1 || clips[0].Timestamp != "2020-04-17 16:16:36" {
		t.Fatalf("Yi quarantine: %+v %v", clips, err)
	}
	if err := os.RemoveAll(filepath.Join(cam.Dir, ".neonvr_pruned")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(cam.Dir, ".neonvr_pruned")); err != nil {
		t.Fatal(err)
	}
	clips, err = s.scanClipsInView(nil, "", false, true)
	if err != nil || len(clips) != 0 {
		t.Fatal("quarantine listing followed symlink root")
	}
}

func TestAllDatesIncludeQuarantineAndFilterSelectedDate(t *testing.T) {
	s, cam, _, _ := detectionFixture(t)
	dir := filepath.Join(cam.Dir, ".neonvr_pruned")
	if err := os.Mkdir(dir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "2026-09-30_12-00-00_old.mp4"), []byte("clip"), 0640); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.apiDatesHandler(w, httptest.NewRequest(http.MethodGet, "/api/dates?view=all", nil))
	var dates datesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &dates); err != nil || len(dates.Dates) != 2 || dates.Dates[1] != "2026-09-30" {
		t.Fatalf("All missing quarantine-only date: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	s.apiAllHandler(w, httptest.NewRequest(http.MethodGet, "/api/all?date=2026-09-30", nil))
	var response clipsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || len(response.Clips) != 1 || !response.Clips[0].Quarantined || response.Clips[0].Date != "2026-09-30" {
		t.Fatalf("All date filtering: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	s.apiDatesHandler(w, httptest.NewRequest(http.MethodGet, "/api/dates?view=clips", nil))
	if err := json.Unmarshal(w.Body.Bytes(), &dates); err != nil || len(dates.Dates) != 1 {
		t.Fatalf("Clips included quarantine-only date: %s", w.Body.String())
	}
}
