package models

import "time"

type VideoData struct {
	VideoID    *string
	Thumbnail  *Image
	Title      *string
	Author     *string
	Duration   *time.Duration
	VidQuality []string
}

type Image struct {
	Src    string
	Width  uint
	Height uint
}
