package handler

import (
	"context"
	"ytdl_http/models"
)

type Servicer interface {
	GetStatus(videoID string) string
	GetInfo(url string) ([]models.VideoData, error)
	DownloadInfo(videoID string) ([]string, error)
	DownloadAudio(ctx context.Context, id string) error
	DownloadVideo(ctx context.Context, id, qual string) error
}
