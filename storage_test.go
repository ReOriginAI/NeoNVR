package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStorageCleanupPriorityAndProtection(t *testing.T) {
	s, cam, original, result := detectionFixture(t)
	s.detectionConfig.StorageMinFreePercent = 10
	contents, err := os.ReadFile(original)
	if err != nil {
		t.Fatal(err)
	}
	makeClip := func(file string, kinds []string, uncertain bool) string {
		t.Helper()
		path := filepath.Join(cam.Dir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, contents, 0640); err != nil {
			t.Fatal(err)
		}
		old := time.Now().Add(-time.Hour)
		if err := os.Chtimes(path, old, old); err != nil {
			t.Fatal(err)
		}
		if kinds != nil {
			d := result
			d.Model = "yolo26s"
			d.Uncertain = uncertain
			d.Intervals = []DetectionInterval{}
			for _, kind := range kinds {
				d.Intervals = append(d.Intervals, DetectionInterval{Kind: kind, Labels: []string{kind}, Start: 0, End: 1, Confidence: .9})
			}
			info, _ := os.Stat(path)
			if err := s.embedDetection(context.Background(), path, info, d); err != nil {
				t.Fatal(err)
			}
		}
		return path
	}
	quarantine := makeClip(".neonvr_pruned/empty.mp4", nil, false)
	car := makeClip("car.mp4", []string{"moving_car"}, false)
	older := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(car, older, older); err != nil {
		t.Fatal(err)
	}
	bike := makeClip("bike.mp4", []string{"bike"}, false)
	person := makeClip("person.mp4", []string{"person", "moving_car"}, false)
	animal := makeClip("animal.mp4", []string{"animal"}, false)
	uncertain := makeClip("uncertain.mp4", []string{"moving_car"}, true)
	favorite := makeClip("favorite.mp4", []string{"moving_car"}, false)
	qFavorite := makeClip(".neonvr_pruned/favorite-quarantine.mp4", nil, false)
	for _, file := range []string{"favorite.mp4", "favorite-quarantine.mp4"} {
		if err := s.meta.SetFavorite(clipKey(cam.Name, file), true); err != nil {
			t.Fatal(err)
		}
	}
	// Older vehicles must still outlive quarantine. Stop as soon as reserve recovers.
	s.runStorageCleanup(func(string) (float64, error) {
		if _, err := os.Stat(quarantine); os.IsNotExist(err) {
			return 12, nil
		}
		return 5, nil
	})
	if _, err := os.Stat(quarantine); !os.IsNotExist(err) {
		t.Fatal("quarantine not deleted first")
	}
	if _, err := os.Stat(car); err != nil {
		t.Fatal("vehicle deleted after reserve restored")
	}
	// Constant pressure exhausts eligible footage but keeps all protected clips.
	s.runStorageCleanup(func(string) (float64, error) { return 5, nil })
	for _, path := range []string{car, bike} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("eligible clip retained: %s", path)
		}
	}
	for _, path := range []string{original, person, animal, uncertain, favorite, qFavorite} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("protected clip removed: %s", path)
		}
	}
}

func TestStorageCleanupDisabledOrHealthy(t *testing.T) {
	s, cam, path, _ := detectionFixture(t)
	dir := filepath.Join(cam.Dir, ".neonvr_pruned")
	if err := os.MkdirAll(dir, 0750); err != nil {
		t.Fatal(err)
	}
	q := filepath.Join(dir, "empty.mp4")
	if err := os.Rename(path, q); err != nil {
		t.Fatal(err)
	}
	s.detectionConfig.StorageMinFreePercent = 0
	s.runStorageCleanup(func(string) (float64, error) { return 5, nil })
	s.detectionConfig.StorageMinFreePercent = 10
	s.runStorageCleanup(func(string) (float64, error) { return 10, nil })
	if _, err := os.Stat(q); err != nil {
		t.Fatal("healthy disk lost quarantine")
	}
}

func TestUnconfirmedQuarantineThreshold(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		free                     float64
		favorite, wantQuarantine bool
	}{
		{"exactly half used", 50, false, false},
		{"below half used", 51, false, false},
		{"above half used", 49, false, true},
		{"favorite preserved", 49, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, cam, path, result := detectionFixture(t)
			// The new quarantine rule also works with space deletion disabled and with
			// the separate no-observation policy set to delete.
			s.detectionConfig.StorageMinFreePercent = 0
			s.detectionConfig.DetectionPrune = "delete"
			result.Uncertain = true
			result.Policy = s.detectionPolicy()
			info, _ := os.Stat(path)
			if err := s.embedDetection(context.Background(), path, info, result); err != nil {
				t.Fatal(err)
			}
			if tc.favorite {
				if err := s.meta.SetFavorite(clipKey(cam.Name, filepath.Base(path)), true); err != nil {
					t.Fatal(err)
				}
			}
			s.runStorageCleanup(func(string) (float64, error) { return tc.free, nil })
			_, err := os.Stat(path)
			if os.IsNotExist(err) != tc.wantQuarantine {
				t.Fatalf("quarantined=%t, want %t", os.IsNotExist(err), tc.wantQuarantine)
			}
			if tc.wantQuarantine {
				q := filepath.Join(cam.Dir, ".neonvr_pruned", filepath.Base(path))
				embedded, err := s.embeddedDetection(context.Background(), q)
				if err != nil || !unconfirmedOnly(embedded) {
					t.Fatalf("analysis lost: %+v %v", embedded, err)
				}
			}
		})
	}
}
