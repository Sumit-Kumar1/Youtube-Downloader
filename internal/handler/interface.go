package handler

import (
	"ytdl_http/internal/models"
)

type Servicer interface {
	GetInfo(url string) ([]models.Video, error)
	DownloadInfo(videoID string) ([]string, error)
	Download(id, qual, audioOnly string) error
}
