package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

const cameraMotionMetadataTag = "yi_motion"

var errUnsupportedCameraMotion = errors.New("unsupported camera metadata format; compact schema/data arrays are required")

type motionSample struct {
	Time        float64 `json:"time"`
	UnixMS      int64   `json:"unix_ms"`
	Scene       int     `json:"scene"`
	Move        int     `json:"move"`
	MovingLevel int     `json:"moving_level"`
	BinImgRatio float64 `json:"bin_img_ratio"`
	MovingTh    int     `json:"moving_th"`
}

type cameraMotion struct {
	Duration        float64             `json:"duration"`
	CreationTime    string              `json:"creation_time"`
	Schema          []string            `json:"schema"`
	Data            [][]json.RawMessage `json:"data"`
	Samples         []motionSample      `json:"samples"`
	InvalidRows     int                 `json:"invalid_rows"`
	OutsideClipRows int                 `json:"outside_clip_rows"`
}

func motionNumber(row []json.RawMessage, index int) (float64, bool) {
	var value float64
	if index < 0 || index >= len(row) || string(row[index]) == "null" {
		return 0, false
	}
	if err := json.Unmarshal(row[index], &value); err != nil || !finite(value) {
		return 0, false
	}
	return value, true
}

// Compact rows are measurements at absolute UTC times, not START/STOP events.
// Use the MP4 creation_time, never local filenames or the first retained row,
// as the playback origin. Keep the original header and raw rows intact.
func parseCameraMotion(tag string, duration float64, creationTime string) (*cameraMotion, error) {
	var raw struct {
		Schema json.RawMessage `json:"schema"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(tag), &raw); err != nil {
		return nil, fmt.Errorf("invalid yi_motion JSON: %w", err)
	}
	var schema []string
	if json.Unmarshal(raw.Schema, &schema) != nil || len(schema) == 0 {
		return nil, errUnsupportedCameraMotion
	}
	var rows [][]json.RawMessage
	if len(raw.Data) == 0 || string(raw.Data) == "null" {
		return nil, errors.New("yi_motion data must be an array")
	}
	if err := json.Unmarshal(raw.Data, &rows); err != nil {
		return nil, fmt.Errorf("invalid yi_motion data: %w", err)
	}
	columns := map[string]int{}
	for i, name := range schema {
		if _, exists := columns[name]; exists {
			return nil, fmt.Errorf("duplicate yi_motion column: %s", name)
		}
		columns[name] = i
	}
	required := []string{"unix_ms", "Scene", "Move", "MovingLevel", "BinImgRatio", "MovingTh"}
	for _, name := range required {
		if _, exists := columns[name]; !exists {
			return nil, fmt.Errorf("yi_motion column is missing: %s", name)
		}
	}
	if !finite(duration) || duration <= 0 {
		return nil, errors.New("yi_motion has no valid video duration")
	}
	result := &cameraMotion{Duration: duration, CreationTime: creationTime, Schema: schema, Data: rows, Samples: []motionSample{}}
	// Empty data reports no retained readings even if the MP4 has no time origin.
	if len(rows) == 0 {
		return result, nil
	}
	start, err := time.Parse(time.RFC3339Nano, creationTime)
	if err != nil {
		return nil, errors.New("camera samples require a valid UTC video creation_time")
	}
	result.CreationTime = start.UTC().Format(time.RFC3339Nano)
	for _, row := range rows {
		values := make([]float64, len(required))
		valid := true
		for i, name := range required {
			value, ok := motionNumber(row, columns[name])
			values[i] = value
			valid = valid && ok
		}
		// Validate documented vendor ranges; unavailable readings remain gaps.
		bounds := [][2]float64{{0, 9007199254740991}, {0, 2}, {0, 4}, {0, 3}, {0, 100}, {10, 31}}
		for i, value := range values {
			valid = valid && value >= bounds[i][0] && value <= bounds[i][1] && (i == 4 || value == math.Trunc(value))
		}
		if !valid {
			result.InvalidRows++
			continue
		}
		offset := (values[0]-float64(start.UnixMilli()))/1000 - float64(start.Nanosecond()%1000000)/1e9
		if offset < 0 || offset >= duration {
			result.OutsideClipRows++
			continue
		}
		result.Samples = append(result.Samples, motionSample{offset, int64(values[0]), int(values[1]), int(values[2]), int(values[3]), values[4], int(values[5])})
	}
	sort.SliceStable(result.Samples, func(i, j int) bool { return result.Samples[i].Time < result.Samples[j].Time })
	return result, nil
}

func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
		"format=duration:format_tags="+cameraMotionMetadataTag+",creation_time:stream=codec_type:stream_tags=creation_time", "-of", "json", path).Output()
	if err != nil {
		return nil, fmt.Errorf("probe camera motion: %w", err)
	}
	var probe struct {
		Streams []struct {
			CodecType string            `json:"codec_type"`
			Tags      map[string]string `json:"tags"`
		} `json:"streams"`
		Format struct {
			Duration string            `json:"duration"`
			Tags     map[string]string `json:"tags"`
		} `json:"format"`
	}
	if err := json.Unmarshal(out, &probe); err != nil {
		return nil, err
	}
	tag := probe.Format.Tags[cameraMotionMetadataTag]
	if tag == "" {
		return nil, nil
	}
	duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
	creationTime := probe.Format.Tags["creation_time"]
	if _, err := time.Parse(time.RFC3339Nano, creationTime); err != nil {
		for _, stream := range probe.Streams {
			if stream.CodecType != "video" {
				continue
			}
			if _, err := time.Parse(time.RFC3339Nano, stream.Tags["creation_time"]); err == nil {
				creationTime = stream.Tags["creation_time"]
				break
			}
		}
	}
	return parseCameraMotion(tag, duration, creationTime)
}

func (s *Server) apiCameraMotionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
		return
	}
	path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	motion, err := s.readCameraMotion(r.Context(), path)
	if errors.Is(err, errUnsupportedCameraMotion) {
		writeJSON(w, map[string]any{"motion": nil, "unsupported": true, "message": err.Error()})
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, struct {
		Motion *cameraMotion `json:"motion"`
	}{motion})
}

func (s *Server) runPreferredThumbnail(videoPath, tmpPath string) error {
	motion, err := s.readCameraMotion(s.workContext(), videoPath)
	if err == nil && motion != nil {
		for _, sample := range motion.Samples {
			// These are observed readings, not inferred motion starts.
			seek := sample.Time
			err = s.runFFmpegThumbnail(videoPath, tmpPath, strconv.FormatFloat(seek, 'f', 3, 64))
			if err == nil {
				var frame frameSignature
				frame, err = frameLumaSignature(tmpPath)
				if err == nil && !frame.usable {
					err = errors.New("camera motion frame is blank")
				}
			}
			if err == nil {
				log.Printf("thumbnail selected %.3fs for %s (camera sample)", seek, videoPath)
				return nil
			}
			_ = os.Remove(tmpPath)
		}
	}
	if err != nil {
		log.Printf("camera motion thumbnail unavailable for %s, using scene selection: %v", videoPath, err)
	}
	return s.runInterestingThumbnail(videoPath, tmpPath)
}
