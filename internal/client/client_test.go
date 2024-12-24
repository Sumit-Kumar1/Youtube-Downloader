package client

import (
	"errors"
	"testing"
	"time"
	"ytdl_http/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
)

var (
	errVid = errors.New("vid error")
)

func TestClient_GetVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockYtdlr := NewMockYtdlr(ctrl)

	defer ctrl.Finish()

	url := "dummy.Youtube.com"
	id := uuid.NewString()
	title := "dummy Vid"
	dur, err := time.ParseDuration("3m")
	if err != nil {
		t.Errorf("error during parsing time: %v", err)
	}

	vid := youtube.Video{
		ID:         id,
		Title:      title,
		Author:     title,
		Duration:   dur,
		Thumbnails: []youtube.Thumbnail{{URL: url, Width: 100, Height: 100}},
	}

	tests := []struct {
		name    string
		url     string
		mock    func()
		want    *models.Video
		wantErr error
	}{
		{name: "nil case", url: url, mock: func() { mockYtdlr.EXPECT().GetVideo(url).Return(nil, nil) }},
		{name: "valid case", url: url, mock: func() { mockYtdlr.EXPECT().GetVideo(url).Return(&vid, nil) },
			want: &models.Video{ID: vid.ID, Title: vid.Title, Duration: dur, Author: vid.Author,
				Thumbnail: models.Image{URL: url, Width: 100, Height: 100}}},
		{name: "err case", url: url, mock: func() { mockYtdlr.EXPECT().GetVideo(url).Return(nil, errVid) }, wantErr: errVid},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{Ytdl: mockYtdlr}

			tt.mock()

			got, err := c.GetVideo(tt.url)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestClient_GetPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockYtdlr := NewMockYtdlr(ctrl)

	defer ctrl.Finish()

	url := "dummy.Youtube.com"
	id := uuid.NewString()
	title := "dummy Vid"

	ytPl := youtube.Playlist{
		ID: id, Title: title, Description: title, Author: title,
		Videos: []*youtube.PlaylistEntry{nil},
	}

	tests := []struct {
		name    string
		url     string
		mock    func()
		want    *models.Playlist
		wantErr error
	}{
		{name: "err case", url: url, mock: func() { mockYtdlr.EXPECT().GetPlaylist(url).Return(nil, errVid) }, want: nil, wantErr: errVid},
		{name: "nil case", url: url, mock: func() { mockYtdlr.EXPECT().GetPlaylist(url).Return(nil, nil) }, want: nil, wantErr: nil},
		{name: "empty case", url: url, mock: func() { mockYtdlr.EXPECT().GetPlaylist(url).Return(&ytPl, nil) }, want: &models.Playlist{
			ID: id, Title: title, Author: title, Description: title}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{Ytdl: mockYtdlr}

			tt.mock()

			got, err := c.GetPlaylist(tt.url)

			assert.Equalf(t, tt.wantErr, err, "Test[%d] Failed - %s", i, tt.name)
			assert.Equalf(t, tt.want, got, "Test[%d] Failed - %s", i, tt.name)
		})
	}
}
