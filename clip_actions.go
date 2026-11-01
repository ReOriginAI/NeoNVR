package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func isQuarantinedFile(file string) bool {
	return strings.HasPrefix(filepath.ToSlash(filepath.Clean(filepath.FromSlash(file))), ".neonvr_pruned/")
}

func (s *Server) apiClipLocationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
		return
	}
	var req struct {
		Camera     string `json:"camera"`
		File       string `json:"file"`
		Quarantine bool   `json:"quarantine"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	file, err := s.moveClip(req.Camera, req.File, req.Quarantine)
	if err != nil {
		writeJSONError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, struct {
		OK   bool   `json:"ok"`
		File string `json:"file"`
	}{true, file})
}

func (s *Server) downloadHandler(w http.ResponseWriter, r *http.Request) {
	path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(path)}))
	http.ServeFile(w, r, path)
}

func (m *MetaStore) moveClip(from, to string, recovered bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	oldFrom, hasFrom := m.data.Clips[from]
	oldTo, hasTo := m.data.Clips[to]
	meta := oldFrom
	meta.Favorite = meta.Favorite || oldTo.Favorite
	meta.Recovered = recovered
	meta.UpdatedAt = time.Now().Format(time.RFC3339)
	delete(m.data.Clips, from)
	m.data.Clips[to] = meta
	if err := m.saveLocked(); err != nil {
		delete(m.data.Clips, to)
		if hasFrom {
			m.data.Clips[from] = oldFrom
		}
		if hasTo {
			m.data.Clips[to] = oldTo
		}
		return err
	}
	return nil
}

// Explicit moves preserve media bytes and timestamps, never overwrite a clip,
// and transfer its favorites and review state to the new relative path.
func (s *Server) moveClip(camera, file string, quarantine bool) (string, error) {
	s.clipMutationMu.Lock()
	defer s.clipMutationMu.Unlock()
	path, err := s.resolveClipPath(camera, file)
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("move requires a regular recording")
	}
	file = filepath.ToSlash(filepath.Clean(filepath.FromSlash(file)))
	if isQuarantinedFile(file) == quarantine {
		return "", fmt.Errorf("clip is already in the requested location")
	}
	cameras, err := s.listCameras()
	if err != nil {
		return "", err
	}
	var root string
	for _, cam := range cameras {
		if cam.Name == camera {
			root = cam.Dir
			break
		}
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	targetFile := strings.TrimPrefix(file, ".neonvr_pruned/")
	if quarantine {
		targetFile = ".neonvr_pruned/" + file
	}
	// Refuse symlink directories before creating or writing any destination.
	dir := root
	for _, part := range strings.Split(filepath.ToSlash(filepath.Dir(targetFile)), "/") {
		if part == "." {
			continue
		}
		dir = filepath.Join(dir, part)
		if err := os.Mkdir(dir, 0750); err != nil && !os.IsExist(err) {
			return "", err
		}
		entry, err := os.Lstat(dir)
		if err != nil || !entry.IsDir() {
			return "", fmt.Errorf("destination directory must not be a symlink")
		}
	}
	thumb, _ := s.thumbnailPath(path)
	target := filepath.Join(root, filepath.FromSlash(targetFile))
	if err := os.Link(path, target); err != nil {
		return "", fmt.Errorf("cannot move clip (destination may already exist): %w", err)
	}
	if err := s.meta.moveClip(clipKey(camera, file), clipKey(camera, targetFile), !quarantine); err != nil {
		_ = os.Remove(target)
		return "", err
	}
	if err := os.Remove(path); err != nil {
		// Retain the original footage if unlinking fails; remove the extra link.
		_ = os.Remove(target)
		_ = s.meta.moveClip(clipKey(camera, targetFile), clipKey(camera, file), false)
		return "", err
	}
	if thumb != "" {
		_ = os.Remove(thumb)
	}
	_ = os.Remove(s.detectionCachePath(path))
	s.invalidateSizeCache()
	return targetFile, nil
}
