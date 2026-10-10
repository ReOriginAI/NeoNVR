package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func fixture(t *testing.T) (Config, string, *configManager) {
	t.Helper()
	c := defaultConfig()
	c.Root = t.TempDir()
	path := filepath.Join(t.TempDir(), "neonvr.json")
	if err := atomicWriteConfig(path, configBytes(c)); err != nil {
		t.Fatal(err)
	}
	manager, err := newConfigManager(path, c)
	if err != nil {
		t.Fatal(err)
	}
	return c, path, manager
}

func postConfig(t *testing.T, m *configManager, c Config, revision, token, origin string, restart bool) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"config": c, "revision": revision, "restart": restart})
	req := httptest.NewRequest(http.MethodPost, "http://neonvr.test/api/config", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-NeoNVR-Token", token)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	m.handler(w, req)
	return w
}

func TestSaveConfigValidationConflictAndRollback(t *testing.T) {
	c, path, m := fixture(t)
	original, _ := os.ReadFile(path)
	revision := configRevision(original)
	c.ThumbWorkers = 0
	w := postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
	if w.Code != 400 {
		t.Fatalf("invalid config status %d: %s", w.Code, w.Body.String())
	}
	raw, _ := os.ReadFile(path)
	if !bytes.Equal(raw, original) {
		t.Fatal("invalid settings changed the file")
	}
	c.ThumbWorkers = 4
	w = postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
	if w.Code != 200 {
		t.Fatalf("save status %d: %s", w.Code, w.Body.String())
	}
	previous, _ := os.ReadFile(path + ".previous")
	if !bytes.Equal(previous, original) {
		t.Fatal("previous config not preserved")
	}
	saved, _ := os.ReadFile(path)
	loaded, err := decodeConfig(saved)
	if err != nil || loaded.ThumbWorkers != 4 {
		t.Fatal("config not saved")
	}
	if m.active.ThumbWorkers != 2 {
		t.Fatal("save applied config without restart")
	}
	if postConfig(t, m, c, revision, m.token, "", false).Code != 409 {
		t.Fatal("stale edit did not conflict")
	}
}

func TestConfigRejectsCrossOriginAndMissingToken(t *testing.T) {
	c, path, m := fixture(t)
	raw, _ := os.ReadFile(path)
	rev := configRevision(raw)
	for _, tc := range []struct{ token, origin string }{{"", "http://neonvr.test"}, {m.token, "http://evil.test"}, {m.token, "https://neonvr.test"}} {
		if w := postConfig(t, m, c, rev, tc.token, tc.origin, false); w.Code != 403 {
			t.Fatalf("expected 403 got %d", w.Code)
		}
	}
}

func TestConcurrentConfigSaves(t *testing.T) {
	c, path, m := fixture(t)
	raw, _ := os.ReadFile(path)
	rev := configRevision(raw)
	var wg sync.WaitGroup
	statuses := make(chan int, 2)
	for i := 3; i < 5; i++ {
		wg.Add(1)
		go func(workers int) {
			defer wg.Done()
			copy := c
			copy.ThumbWorkers = workers
			statuses <- postConfig(t, m, copy, rev, m.token, "", false).Code
		}(i)
	}
	wg.Wait()
	close(statuses)
	count := map[int]int{}
	for status := range statuses {
		count[status]++
	}
	if count[200] != 1 || count[409] != 1 {
		t.Fatalf("unexpected concurrent results %v", count)
	}
}

func TestStartupFileAndLegacyOverrides(t *testing.T) {
	c, path, _ := fixture(t)
	c.ThumbWorkers = 7
	c.ExtraRoots = []string{t.TempDir()}
	os.WriteFile(path, configBytes(c), 0640)
	loaded, _, err := loadStartupConfig([]string{"-config", path})
	if err != nil || loaded.ThumbWorkers != 7 || len(loaded.ExtraRoots) != 1 {
		t.Fatalf("startup: %+v %v", loaded, err)
	}
	replacement := t.TempDir()
	loaded, _, err = loadStartupConfig([]string{"-config", path, "-thumb-workers", "3", "-extra-root", replacement})
	if err != nil || loaded.ThumbWorkers != 3 || len(loaded.ExtraRoots) != 1 || loaded.ExtraRoots[0] != replacement {
		t.Fatalf("override: %+v %v", loaded, err)
	}
	loaded, _, err = loadStartupConfig([]string{"-config", path})
	if err != nil || loaded.ThumbWorkers != 3 {
		t.Fatal("legacy overrides were not persisted")
	}
}

func TestConfigStrictJSONAndLimits(t *testing.T) {
	for _, data := range []string{"null", "{} {}", "{\"typo\":1}", "[]"} {
		if _, err := decodeConfig([]byte(data)); err == nil {
			t.Fatalf("accepted invalid JSON %s", data)
		}
	}
	c, _, _ := fixture(t)
	cases := []func(*Config){func(c *Config) { c.AutoDeleteDays = -1 }, func(c *Config) { c.SizeRefresh = "0" }, func(c *Config) { c.ThumbMaintenanceEvery = "-1h" }, func(c *Config) { c.ThumbWidth = 1 }, func(c *Config) { c.Addr = ":0" }, func(c *Config) { c.Root = "relative" }}
	for _, change := range cases {
		copy := c
		change(&copy)
		if err := copy.validate(); err == nil {
			t.Fatal("invalid config accepted")
		}
	}
	c.ThumbMaintenanceEvery = "0"
	if err := c.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAtomicConfigWriteRejectsSymlink(t *testing.T) {
	target := filepath.Join(t.TempDir(), "target")
	os.WriteFile(target, []byte("original"), 0600)
	link := filepath.Join(t.TempDir(), "config.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := atomicWriteConfig(link, []byte("new")); err == nil {
		t.Fatal("symlink config accepted")
	}
	raw, _ := os.ReadFile(target)
	if string(raw) != "original" {
		t.Fatal("symlink target changed")
	}
}

func TestBackgroundLoopsStopForRestart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{runCtx: ctx, thumbMaintenanceEvery: time.Minute, autoDeleteEvery: time.Minute}
	done := make(chan struct{}, 2)
	go func() { server.thumbnailMaintenanceLoop(); done <- struct{}{} }()
	go func() { server.autoDeleteLoop(); done <- struct{}{} }()
	cancel()
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("background loop did not stop")
		}
	}
}
