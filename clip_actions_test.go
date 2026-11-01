package main

import (
	"bytes"
	"context"
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestManualQuarantineRecoverAndDownload(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Base(path)
	if err := s.meta.SetFavorite(clipKey(cam.Name, file), true); err != nil {
		t.Fatal(err)
	}
	move := func(file string, q bool) string {
		t.Helper()
		body, _ := json.Marshal(map[string]any{"camera": cam.Name, "file": file, "quarantine": q})
		w := httptest.NewRecorder()
		s.apiClipLocationHandler(w, httptest.NewRequest(http.MethodPost, "/api/clip-location", bytes.NewReader(body)))
		var response struct {
			File string `json:"file"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 {
			t.Fatalf("move API: %d %s", w.Code, w.Body.String())
		}
		return response.File
	}
	qfile := move(file, true)
	if qfile != ".neonvr_pruned/"+file {
		t.Fatal(qfile)
	}
	if !s.meta.Get(clipKey(cam.Name, qfile)).Favorite || s.meta.Get(clipKey(cam.Name, file)).Favorite {
		t.Fatal("favorite did not move")
	}
	// All lists both locations and explicitly marks quarantine.
	if err := os.WriteFile(filepath.Join(cam.Dir, "other.mp4"), before, 0640); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.apiAllHandler(w, httptest.NewRequest(http.MethodGet, "/api/all", nil))
	var response clipsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || len(response.Clips) != 2 {
		t.Fatalf("All API: %s", w.Body.String())
	}
	var quarantined Clip
	for _, clip := range response.Clips {
		if clip.Quarantined {
			quarantined = clip
		}
	}
	if quarantined.File != qfile {
		t.Fatal("All missing quarantine marker")
	}
	w = httptest.NewRecorder()
	s.downloadHandler(w, httptest.NewRequest(http.MethodGet, "/download?camera="+cam.Name+"&file="+qfile, nil))
	disposition, params, err := mime.ParseMediaType(w.Header().Get("Content-Disposition"))
	if err != nil || disposition != "attachment" || params["filename"] != file || !bytes.Equal(w.Body.Bytes(), before) {
		t.Fatal("download changed bytes or missing attachment header")
	}
	recovered := move(qfile, false)
	if recovered != file {
		t.Fatal(recovered)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("move changed recording")
	}
	if err := s.meta.SetFavorite(clipKey(cam.Name, file), false); err != nil {
		t.Fatal(err)
	}
	if !s.meta.Get(clipKey(cam.Name, file)).Recovered {
		t.Fatal("unfavorite lost recovery decision")
	}
	if err := s.pruneDetectionClip(context.Background(), cam, file, result); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("automatic no-target pruning undid recovery")
	}
}

func TestManualMovesRefuseCollisionAndExternalSymlink(t *testing.T) {
	for _, symlink := range []bool{false, true} {
		t.Run(map[bool]string{false: "collision", true: "symlink"}[symlink], func(t *testing.T) {
			s, cam, path, _ := detectionFixture(t)
			qdir := filepath.Join(cam.Dir, ".neonvr_pruned")
			if symlink {
				if err := os.Symlink(t.TempDir(), qdir); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.Mkdir(qdir, 0750); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(qdir, filepath.Base(path)), []byte("existing"), 0640); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := s.moveClip(cam.Name, filepath.Base(path), true); err == nil {
				t.Fatal("unsafe move allowed")
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatal("source lost")
			}
			if !symlink {
				data, _ := os.ReadFile(filepath.Join(qdir, filepath.Base(path)))
				if string(data) != "existing" {
					t.Fatal("destination overwritten")
				}
			}
		})
	}
}

func TestRecoveredUnconfirmedStaysInLibrary(t *testing.T) {
	s, cam, path, result := detectionFixture(t)
	result.Uncertain = true
	info, _ := os.Stat(path)
	if err := s.embedDetection(context.Background(), path, info, result); err != nil {
		t.Fatal(err)
	}
	qfile, err := s.moveClip(cam.Name, filepath.Base(path), true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.moveClip(cam.Name, qfile, false); err != nil {
		t.Fatal(err)
	}
	s.runStorageCleanup(func(string) (float64, error) { return 49, nil })
	if _, err = os.Stat(path); err != nil {
		t.Fatal("space policy undid recovery")
	}
}
