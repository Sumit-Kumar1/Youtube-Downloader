package service

import (
	"ytdl_http/internal/models"
)

type Service struct {
	YtClient YtClient
}

func New(yc YtClient) *Service {
	return &Service{YtClient: yc}
}

func (s *Service) GetInfo(url string) ([]models.Video, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}

	if isPlaylistURL(url) {
		return s.getPlaylistData(url)
	}

	return s.getVideoData(url)
}

func (s *Service) DownloadInfo(videoID string) ([]string, error) {
	return s.YtClient.GetDownloadInfo(videoID)
}

func (s *Service) Download(id, qual, audioOnly string) error {
	switch audioOnly {
	case "":
		if !isFFMpegInstalled() {
			return models.ErrNotFound("'ffmpeg' executable")
		}

		if err := s.YtClient.DownloadVideo(id, qual); err != nil {
			return err
		}

	case "true":
		if err := s.YtClient.DownloadAudio(id); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) getPlaylistData(url string) ([]models.Video, error) {
	pl, err := s.YtClient.GetPlaylist(url)
	if err != nil {
		return nil, err
	}

	if pl == nil {
		return nil, nil
	}

	return pl.Videos, nil
}

func (s *Service) getVideoData(url string) ([]models.Video, error) {
	var vids []models.Video

	vid, err := s.YtClient.GetVideo(url)
	if err != nil {
		return nil, err
	}

	if vid == nil {
		return nil, nil
	}

	vids = append(vids, *vid)

	return vids, nil
}
