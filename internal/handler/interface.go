package handler

import (
	"context"
	"ytdl_http/internal/models"
)

type Servicer interface {
	GetInfo(ctx context.Context, url string) ([]models.Video, error)
	DownloadInfo(ctx context.Context, videoID string) ([]string, error)
	Download(ctx context.Context, id, qual, audioOnly string) error
}
