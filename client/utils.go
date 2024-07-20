package client

import (
	"ytdl_http/models"

	"github.com/kkdai/youtube/v2"
)

func getThumbnail(tbs youtube.Thumbnails) *models.Image {
	tb := tbs[len(tbs)/2]

	img := models.Image{
		Url:    tb.URL,
		Height: tb.Height,
		Width:  tb.Width,
	}

	return &img
}
