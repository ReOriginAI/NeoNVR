package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// Use the space available to this process, rather than root's reserved blocks.
func diskFreePercent(path string) (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	if stat.Blocks == 0 {
		return 100, nil
	}
	return 100 * float64(stat.Bavail) / float64(stat.Blocks), nil
}

func (s *Server) storageCleanupLoop() {
	if !s.waitWork(15 * time.Second) {
		return
	}
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		s.runStorageCleanup(diskFreePercent)
		select {
		case <-s.workContext().Done():
			return
		case <-ticker.C:
		}
	}
}

type storageCandidate struct {
	cam        Camera
	file, path string
	info       os.FileInfo
	quarantine bool
}

func vehicleOnly(result *DetectionResult) bool {
	if result == nil || result.validate() != nil || result.Uncertain || len(result.Intervals) == 0 {
		return false
	}
	for _, interval := range result.Intervals {
		if interval.Kind != "moving_car" && interval.Kind != "bike" {
			return false
		}
	}
	return true
}

// Check actual free space after every deletion. Sources on a shared filesystem
// therefore share one reserve, and separate mounts are handled independently.
func (s *Server) runStorageCleanup(freePercent func(string) (float64, error)) {
	s.quarantineUnconfirmedObservations(freePercent)
	reserve := s.detectionConfig.StorageMinFreePercent
	if reserve <= 0 {
		return
	}
	cameras, err := s.listCameras()
	if err != nil {
		log.Printf("storage cleanup: %v", err)
		return
	}
	low := func(path string) bool {
		free, err := freePercent(path)
		if err != nil {
			log.Printf("storage cleanup: free space for %s: %v", path, err)
			return false
		}
		return free < reserve
	}
	// Complete the quarantine tier across all sources before considering vehicles.
	for _, quarantine := range []bool{true, false} {
		var candidates []storageCandidate
		for _, cam := range cameras {
			if !low(cam.Dir) {
				continue
			}
			if quarantine {
				dir := filepath.Join(cam.Dir, ".neonvr_pruned")
				// WalkDir does not follow symlink directories, including the root.
				_ = filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
					if s.workContext().Err() != nil {
						return filepath.SkipAll
					}
					if walkErr != nil || entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !strings.EqualFold(filepath.Ext(path), ".mp4") {
						return nil
					}
					info, err := entry.Info()
					if err != nil || !info.Mode().IsRegular() {
						return nil
					}
					rel, err := filepath.Rel(dir, path)
					if err == nil {
						candidates = append(candidates, storageCandidate{cam, filepath.ToSlash(rel), path, info, true})
					}
					return nil
				})
			} else {
				files, err := listCameraFiles(cam)
				if err != nil {
					continue
				}
				minAge, _ := time.ParseDuration(s.detectionConfig.DetectionMinAge)
				for _, file := range files {
					if time.Since(file.Info.ModTime()) < minAge {
						continue
					}
					candidates = append(candidates, storageCandidate{cam, file.RelPath, filepath.Join(cam.Dir, filepath.FromSlash(file.RelPath)), file.Info, false})
				}
			}
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].info.ModTime().Before(candidates[j].info.ModTime()) })
		for _, candidate := range candidates {
			if s.workContext().Err() != nil {
				return
			}
			if !low(candidate.cam.Dir) {
				continue
			}
			s.removeStorageCandidate(candidate)
		}
	}
}

func (s *Server) removeStorageCandidate(candidate storageCandidate) {
	s.clipMutationMu.Lock()
	defer s.clipMutationMu.Unlock()
	if s.meta.Get(clipKey(candidate.cam.Name, candidate.file)).Favorite ||
		(candidate.quarantine && s.meta.Get(clipKey(candidate.cam.Name, ".neonvr_pruned/"+candidate.file)).Favorite) {
		return
	}
	// Reject links and replacements, and re-read embedded analysis under the same
	// lock used for detection remuxing and favorite changes.
	current, err := os.Lstat(candidate.path)
	if err != nil || !current.Mode().IsRegular() || !sameRecording(candidate.info, current) {
		return
	}
	if !candidate.quarantine {
		path, err := s.resolveClipPath(candidate.cam.Name, candidate.file)
		if err != nil || path != candidate.path {
			return
		}
		result, err := s.embeddedDetection(s.workContext(), path)
		if err != nil || !vehicleOnly(result) {
			return
		}
		current, err = os.Lstat(path)
		if err != nil || !sameRecording(candidate.info, current) {
			return
		}
	}
	if s.workContext().Err() != nil {
		return
	}
	thumb, _ := s.thumbnailPath(candidate.path)
	if err := os.Remove(candidate.path); err != nil {
		log.Printf("storage cleanup: %s: %v", candidate.path, err)
		return
	}
	if thumb != "" {
		_ = os.Remove(thumb)
	}
	_ = os.Remove(s.detectionCachePath(candidate.path))
	// A quarantined clip can share its original key with a newer recording.
	if !candidate.quarantine {
		_ = s.meta.Remove(clipKey(candidate.cam.Name, candidate.file))
	} else {
		_ = s.meta.Remove(clipKey(candidate.cam.Name, ".neonvr_pruned/"+candidate.file))
	}
	s.invalidateSizeCache()
	log.Printf("storage cleanup: deleted %s (quarantine=%t)", candidate.path, candidate.quarantine)
}

// Moving files into quarantine does not free space. It makes weak-only footage
// available for review, then eligible for the existing low-space deletion tier.
func (s *Server) quarantineUnconfirmedObservations(freePercent func(string) (float64, error)) {
	cameras, err := s.listCameras()
	if err != nil {
		return
	}
	minAge, _ := time.ParseDuration(s.detectionConfig.DetectionMinAge)
	for _, cam := range cameras {
		free, err := freePercent(cam.Dir)
		if err != nil || free >= 50 {
			continue
		}
		files, err := listCameraFiles(cam)
		if err != nil {
			continue
		}
		for _, file := range files {
			if s.workContext().Err() != nil {
				return
			}
			meta := s.meta.Get(clipKey(cam.Name, file.RelPath))
			if time.Since(file.Info.ModTime()) < minAge || meta.Favorite || meta.Recovered {
				continue
			}
			result := s.cachedDetection(file.FullPath, file.Info)
			if result == nil {
				result, _ = s.embeddedDetection(s.workContext(), file.FullPath)
				if result != nil {
					_ = s.cacheDetection(file.FullPath, file.Info, *result)
				}
			}
			if !unconfirmedOnly(result) {
				continue
			}
			if err := s.pruneDetectionClipWithPolicy(s.workContext(), cam, file.RelPath, *result, "quarantine", true); err != nil {
				log.Printf("storage quarantine: %s: %v", file.FullPath, err)
			}
		}
	}
}
