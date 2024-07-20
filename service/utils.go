package service

import (
	"errors"
	"regexp"
	"strings"
)

var (
	playlistRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtube\.com\/playlist\?list=([a-zA-Z0-9_-]+)`)
)

func validateURL(url string) error {
	if !strings.Contains(url, "youtu") || !strings.ContainsAny(url, "\"?&/<%=") {
		return errors.New("not a valid link")
	}

	return nil
}

func isPlaylistURL(url string) bool {
	return playlistRegex.MatchString(url)
}
