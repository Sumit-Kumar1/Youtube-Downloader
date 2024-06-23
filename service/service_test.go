package service

import (
	"errors"
	"testing"
	"ytdl_http/models"

	yt "github.com/kkdai/youtube/v2"
	dlr "github.com/kkdai/youtube/v2/downloader"
	"github.com/stretchr/testify/assert"
)

func newTestResources() *Service {
	c := &yt.Client{MaxRoutines: 8, ChunkSize: yt.Size10Mb}
	d := &dlr.Downloader{OutputDir: "Downloads"}
	s := New(c, d)

	return s
}

func TestServiceGetStatus(t *testing.T) {
	s := newTestResources()
	s.Status = map[string]string{"abcd": "started", "123": "not-started", "000": "downloaded"}
	tests := []struct {
		name    string
		videoID string
		want    string
	}{
		{"download started", "abcd", "started"},
		{"download not started", "123", "not-started"},
		{"video downloaded", "000", "downloaded"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, s.GetStatus(tt.videoID), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestServiceGetInfo(t *testing.T) {
	s := newTestResources()
	tests := []struct {
		name    string
		url     string
		want    []models.VideoData
		wantErr error
	}{
		{"invalid url", "", nil, errors.New("not a valid link")},
		{"video url", vidURL, nil, nil}, //need internet to run these test as no mocks are avail
		{"playlist url", vidURL, nil, nil},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetInfo(tt.url)
			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}
