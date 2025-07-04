package service

import (
	"testing"

	"ytdl_http/internal/models"

	"github.com/stretchr/testify/assert"
)

const (
	plURL  = "https://www.youtube.com/playlist?list=PLZHQObOWTQDN52m7Y21ePrTbvXkPaWVSg"
	vidURL = "https://www.youtube.com/watch?v=wI2sxYte6RA&list=RDCLAK5uy_n41V9iqmjro6caBDuDD1E4eWs5yTb5_OY&index=1"
)

func Test_validateURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		expErr error
	}{
		{name: "playlist url", url: plURL, expErr: nil},
		{name: "video url", url: vidURL, expErr: nil},
		{name: "invalid url", url: "www.lazy.com", expErr: models.ErrInvalid("link")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)

			assert.Equalf(t, tt.expErr, err, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_isPlaylistURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "playlist", url: plURL, want: true},
		{name: "video", url: vidURL, want: false},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPlaylistURL(tt.url)
			assert.Equalf(t, tt.want, got, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
