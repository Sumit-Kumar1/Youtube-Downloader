package service

import (
	"context"

	"ytdl_http/internal/models"
)

type YtClient interface {
	GetVideo(ctx context.Context, url string) (*models.Video, error)
	GetPlaylist(ctx context.Context, url string) (*models.Playlist, error)
	GetDownloadInfo(ctx context.Context, videoID string) ([]string, error)

	DownloadVideo(ctx context.Context, id, qual string) error
	DownloadAudio(ctx context.Context, id string) error
}
