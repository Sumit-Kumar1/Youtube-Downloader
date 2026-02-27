package client

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"ytdl_http/internal/models"

	"github.com/kkdai/youtube/v2"
)

var nameCleanRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]`)

func getThumbnail(tbs youtube.Thumbnails) models.Image {
	if len(tbs) == 0 {
		return models.Image{}
	}

	tb := tbs[len(tbs)/2]

	return models.Image{
		URL:    tb.URL,
		Height: tb.Height,
		Width:  tb.Width,
	}
}

func formatName(title string) string {
	if title == "" {
		return title
	}

	name := nameCleanRegex.ReplaceAllString(title, "")

	return strings.TrimSpace(name)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen]
}

func stream2File(stream io.ReadCloser, fileName string) error {
	destPath := filepath.Clean(filepath.Join(models.DirPath, fileName))

	tempF, err := os.CreateTemp(models.DirPath, "*.tmp")
	if err != nil {
		return err
	}

	tempPath := tempF.Name()

	defer func() {
		if err != nil {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err = io.Copy(tempF, stream); err != nil {
		_ = tempF.Close()

		return err
	}

	if err = tempF.Sync(); err != nil {
		_ = tempF.Close()

		return err
	}

	if err = tempF.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, destPath)
}

func addMetaData(ctx context.Context, vid *youtube.Video, fileName string) error {
	inputPath := filepath.Clean(filepath.Join(models.DirPath, fileName))

	if _, err := os.Stat(inputPath); err != nil {
		return err
	}

	outputPath := inputPath + ".tmp"
	args := prepareMetadataArgs(vid, inputPath, outputPath)

	out, err := exec.CommandContext(ctx, "ffmpeg", args...).CombinedOutput()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "ffmpeg metadata failed",
			slog.String("output", string(out)),
			slog.String("error", err.Error()))

		_ = os.Remove(outputPath)

		return err
	}

	return os.Rename(outputPath, inputPath)
}

func prepareMetadataArgs(vid *youtube.Video, inputFile, outputFile string) []string {
	args := []string{"-y", "-i", inputFile, "-map", "0", "-c:a", "copy"}

	meta := map[string]string{
		"title":       truncate(vid.Title, 100),
		"author":      truncate(vid.Author, 100),
		"description": truncate(vid.Description, 100),
		"year":        strconv.Itoa(vid.PublishDate.Year()),
	}

	for k, v := range meta {
		args = append(args, "-metadata", k+"="+v)
	}

	args = append(args, outputFile)

	return args
}
