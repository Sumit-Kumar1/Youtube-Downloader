package service

import (
	"net/url"
	"os/exec"
	"regexp"
	"sync"

	"ytdl_http/internal/models"
)

var (
	playlistRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtube\.com\/playlist\?list=([a-zA-Z0-9_-]+)`)

	ffmpegOnce      sync.Once
	ffmpegInstalled bool
)

var validYouTubeHosts = map[string]bool{
	"youtube.com":     true,
	"www.youtube.com": true,
	"m.youtube.com":   true,
	"youtu.be":        true,
	"www.youtu.be":    true,
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return models.ErrInvalid("link")
	}

	if !validYouTubeHosts[u.Host] {
		return models.ErrInvalid("link")
	}

	return nil
}

func isPlaylistURL(url string) bool {
	return playlistRegex.MatchString(url)
}

func isFFMpegInstalled() bool {
	ffmpegOnce.Do(func() {
		if _, err := exec.LookPath("ffmpeg"); err == nil {
			ffmpegInstalled = true
		}
	})

	return ffmpegInstalled
}
