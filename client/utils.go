package client

import (
	"regexp"
	"strings"
	"ytdl_http/models"

	"github.com/kkdai/youtube/v2"
)

func getThumbnail(tbs youtube.Thumbnails) *models.Image {
	if len(tbs) == 0 {
		return nil
	}

	tb := tbs[len(tbs)/2]

	img := models.Image{
		URL:    tb.URL,
		Height: tb.Height,
		Width:  tb.Width,
	}

	return &img
}

func formatName(title string) string {
	if title == "" {
		return title
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	name := re.ReplaceAllString(title, "")

	return strings.TrimSpace(name)
}
