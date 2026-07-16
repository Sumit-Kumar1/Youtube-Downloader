package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_FillByName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Player
	}{
		{name: "empty name", input: "", want: Player{}},
		{name: "short name", input: "ab", want: Player{}},
		{name: "mp4 file", input: "video.mp4", want: Player{
			Title: "video.mp4", Path: "/resource/video.mp4",
			IsAudio: false, Type: VideoType,
		}},
		{name: "m4a file", input: "audio.m4a", want: Player{
			Title: "audio.m4a", Path: "/resource/audio.m4a",
			IsAudio: true, Type: AudioType,
		}},
		{name: "mp3 file", input: "song.mp3", want: Player{
			Title: "song.mp3", Path: "/resource/song.mp3",
			IsAudio: true, Type: AudioType,
		}},
		{name: "unknown extension", input: "readme.txt", want: Player{}},
		{name: "partial file", input: "download.part", want: Player{}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Player
			p.FillByName(tt.input)
			assert.Equalf(t, tt.want, p, "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}
