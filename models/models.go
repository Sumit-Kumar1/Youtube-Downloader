package models

import "time"

type Video struct {
	ID         *string
	Thumbnail  *Image
	Title      *string
	Author     *string
	Duration   *time.Duration
	VidQuality []string
}

type Playlist struct {
	ID          string
	Title       string
	Description string
	Author      string
	Videos      []*PlaylistEntry
}

type PlaylistEntry struct {
	ID         string
	Title      string
	Author     string
	Duration   time.Duration
	Thumbnails Image
}

type Image struct {
	Url    string
	Width  uint
	Height uint
}
