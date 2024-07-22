package client

import (
	"regexp"
	"strings"
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

func formatName(title string) string {
	if title == "" {
		return title
	}

	name := strings.Split(title, "|")[0]
	re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)

	name = re.ReplaceAllString(name, "")

	return name
}
