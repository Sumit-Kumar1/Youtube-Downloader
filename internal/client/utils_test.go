package client

import (
	"testing"

	"ytdl_http/internal/models"

	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
)

func Test_formatName(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{name: "name", title: "abcd", want: "abcd"},
		{name: "name - 2", title: "abcd|bcded", want: "abcdbcded"},
		{"name - 3", "abcd&bcd|ed", "abcdbcded"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, formatName(tt.title), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_getThumbnail(t *testing.T) {
	url := "http://www.image.com"
	var (
		w uint = 10
		h uint = 10
	)

	tests := []struct {
		name string
		tbs  youtube.Thumbnails
		want *models.Image
	}{
		{"nil case", nil, nil},
		{"one thumbnail case", youtube.Thumbnails{{URL: url, Width: w, Height: h}},
			&models.Image{URL: url, Width: w, Height: h}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getThumbnail(tt.tbs), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
