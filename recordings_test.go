package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestParseYICamMinuteSecond(t *testing.T) {
	for _, tc := range []struct {
		name           string
		minute, second int
		ok             bool
	}{
		{"16M36S06.mp4", 16, 36, true},
		{"W440M00S60.mp4", 40, 0, true},
		{"00M00S06.mp4", 0, 0, true},
		{"59M59S120.MP4", 59, 59, true},
		{"60M00S06.mp4", 0, 0, false},
		{"16M60S06.mp4", 0, 0, false},
		{"-1M36S06.mp4", 0, 0, false},
		{"16M36S.mp4", 0, 0, false},
		{"recording.mp4", 0, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			minute, second, ok := parseYICamMinuteSecond(tc.name)
			if minute != tc.minute || second != tc.second || ok != tc.ok {
				t.Fatalf("got (%d, %d, %v), want (%d, %d, %v)", minute, second, ok, tc.minute, tc.second, tc.ok)
			}
		})
	}
}

func TestYiRecordingDatesAndSorting(t *testing.T) {
	root := t.TempDir()
	camDir := filepath.Join(root, "yicamBasement")
	hour := "2020Y04M17D16H"
	if err := os.MkdirAll(filepath.Join(camDir, hour), 0755); err != nil {
		t.Fatal(err)
	}
	mtime := time.Date(2026, 10, 22, 17, 23, 0, 0, time.Local)
	for _, rel := range []string{
		hour + "/16M36S06.mp4",
		hour + "/W440M00S60.mp4",
		"2020-04-17_16-30-00_Basement.mp4",
		hour + "/unknown.mp4",
	} {
		path := filepath.Join(camDir, filepath.FromSlash(rel))
		if err := os.WriteFile(path, nil, 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, mtime, mtime); err != nil {
			t.Fatal(err)
		}
	}
	meta, err := NewMetaStore(filepath.Join(t.TempDir(), "meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{root: root, meta: meta}
	dates, err := s.listDates(nil)
	if err != nil || !reflect.DeepEqual(dates, []string{"2026-10-22", "2020-04-17"}) {
		t.Fatalf("dates = %v, error = %v", dates, err)
	}
	clips, err := s.scanClips(nil, "2020-04-17", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(clips) != 3 {
		t.Fatalf("got %d clips, want 3", len(clips))
	}
	for i, want := range []string{"16:40:00", "16:30:00", "16:16:36"} {
		if clips[i].Time != want {
			t.Errorf("clip %d time = %q, want %q", i, clips[i].Time, want)
		}
	}
	want := time.Date(2020, 4, 17, 16, 16, 36, 0, time.Local)
	if clips[2].UnixSeconds != want.Unix() || clips[2].File != hour+"/16M36S06.mp4" {
		t.Fatalf("new Yi clip = %+v", clips[2])
	}
	clips, err = s.scanClips(nil, "2026-10-22", false)
	if err != nil || len(clips) != 1 || clips[0].File != hour+"/unknown.mp4" {
		t.Fatalf("mtime fallback clips = %+v, error = %v", clips, err)
	}
}
