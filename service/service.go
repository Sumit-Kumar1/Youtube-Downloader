package service

import (
	"context"
	"downloader/models"

	yt "github.com/kkdai/youtube/v2"
	dlr "github.com/kkdai/youtube/v2/downloader"
)

type Service struct {
	ytClient   *yt.Client
	downloader *dlr.Downloader
	Status     map[string]string
}

func New(c *yt.Client, d *dlr.Downloader) *Service {
	return &Service{
		ytClient:   c,
		downloader: d,
		Status:     make(map[string]string, 0),
	}
}

// GetStatus get the download status of the video by VideoID
func (s *Service) GetStatus(videoID string) string {
	return s.Status[videoID]
}

// GetInfo provides the video/'s info of the provided url
func (s *Service) GetInfo(url string) ([]models.VideoData, error) {
	isPlaylist, err := validateURL(url)
	if err != nil {
		return nil, err
	}

	// clear status map when we get new request
	for k := range s.Status {
		delete(s.Status, k)
	}

	if isPlaylist {
		return s.getPlaylistInfo(url)
	}

	data := []models.VideoData{}
	d, err := s.getVideoInfo(url)
	if err != nil {
		return nil, err
	}

	data = append(data, *d)

	return data, nil
}

// DownloadInfo return the quality of video for download
func (s *Service) DownloadInfo(videoID string) ([]string, error) {
	vid, err := s.ytClient.GetVideoContext(context.Background(), videoID)
	if err != nil {
		return nil, err
	}

	s.Status[vid.ID] = "Not Started"

	qls := []string{}
	for _, format := range vid.Formats {
		if format.QualityLabel != "" {
			qls = append(qls, format.QualityLabel)
		}
	}

	return qls, nil
}

// Download started download of video in a new thread based on videoID
func (s *Service) DownloadVideo(ctx context.Context, id, qual string) error {
	vid, err := s.ytClient.GetVideoContext(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		s.Status[id] = "Download Started"

		if err := s.downloader.DownloadComposite(context.Background(), vid.Title+".mp4", vid, qual, "", ""); err != nil {
			s.Status[id] = err.Error()
			return
		}

		s.Status[id] = "Download Done"
	}()

	return nil
}

func (s *Service) getVideoInfo(url string) (*models.VideoData, error) {
	vid, err := s.ytClient.GetVideo(url)
	if err != nil {
		return nil, err
	}

	qls := make([]string, 0, len(vid.Formats))
	for _, format := range vid.Formats {
		if format.QualityLabel != "" {
			qls = append(qls, format.QualityLabel)
		}
	}

	tb := vid.Thumbnails[len(vid.Thumbnails)/2]
	data := models.VideoData{
		VideoID:  &vid.ID,
		Title:    &vid.Title,
		Author:   &vid.Author,
		Duration: &vid.Duration,
		Thumbnail: &models.Image{
			Src:    tb.URL,
			Width:  tb.Width,
			Height: tb.Height,
		},
		VidQuality: qls,
	}

	s.Status[vid.ID] = "Not Started"

	return &data, nil
}

func (s *Service) getPlaylistInfo(url string) ([]models.VideoData, error) {
	pl, err := s.ytClient.GetPlaylist(url)
	if err != nil {
		return nil, err
	}

	data := make([]models.VideoData, 0, len(pl.Videos))
	for _, vid := range pl.Videos {
		tb := vid.Thumbnails[len(vid.Thumbnails)/2]
		d := models.VideoData{
			VideoID:  &vid.ID,
			Title:    &vid.Title,
			Author:   &vid.Author,
			Duration: &vid.Duration,
			Thumbnail: &models.Image{
				Src:    tb.URL,
				Height: tb.Height,
				Width:  tb.Width,
			},
		}

		data = append(data, d)
	}

	return data, nil
}
