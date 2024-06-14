package service

import (
	"errors"
	"regexp"
	"strings"
)

var (
	playlistRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtube\.com\/playlist\?list=([a-zA-Z0-9_-]+)`)
)

func validateURL(url string) (bool, error) {
	if !strings.Contains(url, "youtu") || !strings.ContainsAny(url, "\"?&/<%=") {
		return false, errors.New("not a valid link")
	}

	return playlistRegex.MatchString(url), nil
}
