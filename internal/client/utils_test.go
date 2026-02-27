package client

import (
	"testing"
	"time"

	"ytdl_http/internal/models"

	"github.com/kkdai/youtube/v2"
	"github.com/stretchr/testify/assert"
)

func Test_formatName(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{name: "empty title", title: "", want: ""},
		{name: "name", title: "abcd", want: "abcd"},
		{name: "name - 2", title: "abcd|bcded", want: "abcdbcded"},
		{name: "name - 3", title: "abcd&bcd|ed", want: "abcdbcded"},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, formatName(tt.title), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_getThumbnail(t *testing.T) {
	url := "http://www.image.com"

	var (
		w uint = 10
		h uint = 10
	)

	tests := []struct {
		name string
		tbs  youtube.Thumbnails
		want models.Image
	}{
		{name: "nil case", tbs: nil, want: models.Image{}},
		{name: "one thumbnail case", tbs: youtube.Thumbnails{{URL: url, Width: w, Height: h}},
			want: models.Image{URL: url, Width: w, Height: h}},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getThumbnail(tt.tbs), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_truncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{name: "empty string", input: "", maxLen: 100, want: ""},
		{name: "short string", input: "hello", maxLen: 100, want: "hello"},
		{name: "exact length", input: "hello", maxLen: 5, want: "hello"},
		{name: "over length", input: "hello world", maxLen: 5, want: "hello"},
		{name: "200 chars truncated to 100", input: stringOfLen(200), maxLen: 100, want: stringOfLen(100)},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, truncate(tt.input, tt.maxLen), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func Test_prepareMetadataArgs(t *testing.T) {
	vid := &youtube.Video{
		Title:       "Test Title",
		Author:      "Test Author",
		Description: "Test Description",
		PublishDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	args := prepareMetadataArgs(vid, "/input.m4a", "/output.m4a")

	assert.Contains(t, args, "-y")
	assert.Contains(t, args, "-i")
	assert.Contains(t, args, "/input.m4a")
	assert.Contains(t, args, "-map")
	assert.Contains(t, args, "-c:a")
	assert.Contains(t, args, "copy")
	assert.Equal(t, "/output.m4a", args[len(args)-1])

	// Verify metadata key=value pairs are present
	metaStr := ""
	for _, a := range args {
		metaStr += a + " "
	}

	assert.Contains(t, metaStr, "title=Test Title")
	assert.Contains(t, metaStr, "author=Test Author")
	assert.Contains(t, metaStr, "year=2024")
}

func Test_prepareMetadataArgs_Truncation(t *testing.T) {
	longTitle := stringOfLen(200)

	vid := &youtube.Video{
		Title:       longTitle,
		Author:      "short",
		Description: stringOfLen(150),
		PublishDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	args := prepareMetadataArgs(vid, "/in.m4a", "/out.m4a")

	metaStr := ""
	for _, a := range args {
		metaStr += a + " "
	}

	assert.Contains(t, metaStr, "title="+stringOfLen(100))
	assert.Contains(t, metaStr, "author=short")
	assert.Contains(t, metaStr, "description="+stringOfLen(100))
}

func stringOfLen(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}

	return string(b)
}
