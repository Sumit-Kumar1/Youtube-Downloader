package handler

import (
	"context"
	"downloader/models"
)

type Servicer interface {
	GetStatus(videoID string) string
	GetInfo(url string) ([]models.VideoData, error)
	DownloadInfo(videoID string) ([]string, error)
	DownloadVideo(ctx context.Context, id, qual string) error
}
