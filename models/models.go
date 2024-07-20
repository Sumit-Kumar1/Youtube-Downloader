package models

import "time"

type Video struct {
	ID        string
	Author    string
	Title     string
	Thumbnail Image
	Duration  time.Duration
}

type Playlist struct {
	ID          string
	Title       string
	Description string
	Author      string
	Videos      []Video
}

type Image struct {
	Url    string
	Width  uint
	Height uint
}
