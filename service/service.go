package service

import (
	"ytdl_http/models"
)

type Service struct {
	YtClient YtClient
	Status   map[string]string
}

func New(y YtClient) *Service {
	return &Service{
		YtClient: y,
		Status:   make(map[string]string, 0),
	}
}

// GetStatus get the download status of the video by ID
func (s *Service) GetStatus(videoID string) string {
	return s.Status[videoID]
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
	s.Status[id] = "download started"

	switch audioOnly {
	case "":
		if err := s.YtClient.DownloadVideo(id, qual); err != nil {
			s.Status[id] = err.Error()

			return err
		}

	case "true":
		if err := s.YtClient.DownloadAudio(id); err != nil {
			s.Status[id] = err.Error()

			return err
		}
	}

	s.Status[id] = "download complete !!"

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
