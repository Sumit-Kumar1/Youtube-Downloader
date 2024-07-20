package client

import (
	"ytdl_http/models"

	dlr "github.com/kkdai/youtube/v2/downloader"
)

type Client struct {
	Ytdl *dlr.Downloader
}

func New(d *dlr.Downloader) Client {
	return Client{Ytdl: d}
}

func (c Client) GetVideo(url string) (*models.Video, error) {
	ytvid, err := c.Ytdl.GetVideo(url)
	if err != nil {
		return nil, err
	}

	if ytvid == nil {
		return nil, nil
	}

	vid := models.Video{
		ID:        ytvid.ID,
		Author:    ytvid.Author,
		Title:     ytvid.Title,
		Duration:  ytvid.Duration,
		Thumbnail: *getThumbnail(ytvid.Thumbnails),
	}

	return &vid, nil
}

func (c Client) GetPlaylist(url string) (*models.Playlist, error) {
	ytPl, err := c.Ytdl.GetPlaylist(url)
	if err != nil {
		return nil, err
	}

	if ytPl == nil {
		return nil, nil
	}

	playlist := models.Playlist{
		ID:          ytPl.ID,
		Title:       ytPl.Title,
		Description: ytPl.Description,
		Author:      ytPl.Author,
	}

	for _, vid := range ytPl.Videos {
		if vid == nil {
			continue
		}

		v := models.Video{
			ID:        vid.ID,
			Duration:  vid.Duration,
			Title:     vid.Title,
			Author:    vid.Author,
			Thumbnail: *getThumbnail(vid.Thumbnails),
		}

		playlist.Videos = append(playlist.Videos, v)
	}

	return &playlist, nil
}
