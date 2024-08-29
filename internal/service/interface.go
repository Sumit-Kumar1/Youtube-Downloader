package service

import (
	"ytdl_http/internal/models"
)

type YtClient interface {
	GetVideo(url string) (*models.Video, error)
	GetPlaylist(url string) (*models.Playlist, error)
	GetDownloadInfo(videoID string) ([]string, error)
	DownloadVideo(id, quality string) (*models.ClientResponse, error)
	DownloadAudio(id string) (*models.ClientResponse, error)
}
