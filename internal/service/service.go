package service

import (
	"context"

	"ytdl_http/internal/models"
)

type Service struct {
	ytClient YtClient
}

func New(yc YtClient) *Service {
	return &Service{ytClient: yc}
}

func (s *Service) GetInfo(ctx context.Context, url string) ([]models.Video, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}

	if isPlaylistURL(url) {
		return s.getPlaylistData(ctx, url)
	}

	return s.getVideoData(ctx, url)
}

func (s *Service) DownloadInfo(ctx context.Context, videoID string) ([]string, error) {
	return s.ytClient.GetDownloadInfo(ctx, videoID)
}

func (s *Service) Download(ctx context.Context, id, qual string, audioOnly bool) error {
	if !isFFMpegInstalled() {
		return models.ErrNotFound("'ffmpeg' executable")
	}

	if audioOnly {
		return s.ytClient.DownloadAudio(ctx, id)
	}

	return s.ytClient.DownloadVideo(ctx, id, qual)
}

func (s *Service) getPlaylistData(ctx context.Context, url string) ([]models.Video, error) {
	pl, err := s.ytClient.GetPlaylist(ctx, url)
	if err != nil {
		return nil, err
	}

	if pl == nil {
		return nil, nil
	}

	return pl.Videos, nil
}

func (s *Service) getVideoData(ctx context.Context, url string) ([]models.Video, error) {
	vid, err := s.ytClient.GetVideo(ctx, url)
	if err != nil {
		return nil, err
	}

	if vid == nil {
		return nil, nil
	}

	return []models.Video{*vid}, nil
}
