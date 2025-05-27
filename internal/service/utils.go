package service

import (
	"os/exec"
	"regexp"
	"strings"

	"ytdl_http/internal/models"
)

var (
	playlistRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtube\.com\/playlist\?list=([a-zA-Z0-9_-]+)`)
)

func validateURL(url string) error {
	if !strings.Contains(url, "youtu") || !strings.ContainsAny(url, "\"?&/<%=") {
		return models.ErrInvalid("link")
	}

	return nil
}

func isPlaylistURL(url string) bool {
	return playlistRegex.MatchString(url)
}

func isFFMpegInstalled() bool {
	cmd := exec.Command("ffmpeg", "-version")
	if _, err := cmd.CombinedOutput(); err != nil {
		return false
	}

	return true
}
