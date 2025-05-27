package client

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ytdl_http/internal/models"

	"github.com/kkdai/youtube/v2"
)

func getThumbnail(tbs youtube.Thumbnails) *models.Image {
	if len(tbs) == 0 {
		return nil
	}

	tb := tbs[len(tbs)/2]

	img := models.Image{
		URL:    tb.URL,
		Height: tb.Height,
		Width:  tb.Width,
	}

	return &img
}

func formatName(title string) string {
	if title == "" {
		return title
	}

	re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	name := re.ReplaceAllString(title, "")

	return strings.TrimSpace(name)
}

func stream2File(stream io.ReadCloser, fileName string) error {
	tempF, err := os.CreateTemp(os.TempDir(), fileName)
	if err != nil {
		return err
	}

	defer func() { _ = tempF.Close() }()

	if _, err = io.Copy(tempF, stream); err != nil {
		return err
	}

	if err := tempF.Sync(); err != nil {
		return nil
	}

	defer func() { _ = os.RemoveAll(tempF.Name()) }()

	file, err := os.Create(filepath.Clean(models.DirPath + "/" + fileName))
	if err != nil {
		return err
	}

	data, err := os.ReadFile(tempF.Name())
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return file.Close()
}
