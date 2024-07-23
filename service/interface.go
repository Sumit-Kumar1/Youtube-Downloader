package service

import (
	"ytdl_http/models"
)

type YtClient interface {
	GetVideo(url string) (*models.Video, error)
	GetPlaylist(url string) (*models.Playlist, error)
	GetDownloadInfo(videoID string) ([]string, error)
	DownloadVideo(id, qual string) error
	DownloadAudio(id string) error
}
