package client

import (
	"context"
	"ytdl_http/internal/models"

	"github.com/kkdai/youtube/v2"
)

type Client struct {
	Ytdl *youtube.Client
}

func New() Client {
	return Client{Ytdl: &youtube.Client{}}
}

func (c Client) GetVideo(url string) (*models.Video, error) {
	ytVid, err := c.Ytdl.GetVideo(url)
	if err != nil {
		return nil, err
	}

	if ytVid == nil {
		return nil, nil
	}

	vid := models.Video{
		ID:        ytVid.ID,
		Author:    ytVid.Author,
		Title:     ytVid.Title,
		Duration:  ytVid.Duration,
		Thumbnail: *getThumbnail(ytVid.Thumbnails),
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

func (c Client) GetDownloadInfo(videoID string) ([]string, error) {
	vid, err := c.Ytdl.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	fl := vid.Formats

	qls := make([]string, 0, len(fl))
	for i := range fl {
		if fl[i].QualityLabel == "" {
			continue
		}

		qls = append(qls, fl[i].QualityLabel)
	}

	return qls, nil
}

func (c Client) DownloadVideo(id, quality string) (*models.ClientResponse, error) {
	vid, err := c.Ytdl.GetVideo(id)
	if err != nil {
		return nil, err
	}

	name := formatName(vid.Title) + ".mp4"
	formats := vid.Formats.Quality(quality)

	stream, _, err := c.Ytdl.GetStream(vid, &formats[0])
	if err != nil {
		return nil, err
	}

	return &models.ClientResponse{Stream: stream, Filename: name}, nil
}

func (c Client) DownloadAudio(id string) (*models.ClientResponse, error) {
	ctx := context.Background()
	vid, err := c.Ytdl.GetVideoContext(ctx, id)
	if err != nil {
		return nil, err
	}

	title := formatName(vid.Title) + ".m4a"

	audioFormats := vid.Formats.Type("audio")
	audioFormats.Sort()

	stream, _, err := c.Ytdl.GetStreamContext(ctx, vid, &audioFormats[0])
	if err != nil {
		return nil, err
	}

	return &models.ClientResponse{Stream: stream, Filename: title}, nil
}
