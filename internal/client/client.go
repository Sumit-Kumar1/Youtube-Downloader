package client

import (
	"context"
	"fmt"

	"ytdl_http/internal/models"

	dlr "github.com/kkdai/youtube/v2/downloader"
)

type Client struct {
	ytd ytdlr
}

func New() *Client {
	d := dlr.Downloader{OutputDir: models.DirPath}

	return &Client{ytd: &d}
}

func (c *Client) GetVideo(ctx context.Context, url string) (*models.Video, error) {
	ytVid, err := c.ytd.GetVideoContext(ctx, url)
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
		Thumbnail: getThumbnail(ytVid.Thumbnails),
	}

	return &vid, nil
}

func (c *Client) GetPlaylist(ctx context.Context, url string) (*models.Playlist, error) {
	ytPl, err := c.ytd.GetPlaylistContext(ctx, url)
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
			Thumbnail: getThumbnail(vid.Thumbnails),
		}

		playlist.Videos = append(playlist.Videos, v)
	}

	return &playlist, nil
}

func (c *Client) GetDownloadInfo(ctx context.Context, videoID string) ([]string, error) {
	vid, err := c.ytd.GetVideoContext(ctx, videoID)
	if err != nil {
		return nil, err
	}

	qls := make([]string, 0, len(vid.Formats))

	for i := range vid.Formats {
		if vid.Formats[i].QualityLabel != "" {
			qls = append(qls, vid.Formats[i].QualityLabel)
		}
	}

	return qls, nil
}

func (c *Client) DownloadVideo(ctx context.Context, id, qual string) error {
	vid, err := c.ytd.GetVideoContext(ctx, id)
	if err != nil {
		return err
	}

	title := formatName(vid.Title)

	return c.ytd.DownloadComposite(ctx, title+".mp4", vid, qual, "", "")
}

func (c *Client) DownloadAudio(ctx context.Context, id string) error {
	vid, err := c.ytd.GetVideoContext(ctx, id)
	if err != nil {
		return err
	}

	audioFormats := vid.Formats.Type("audio")
	if len(audioFormats) == 0 {
		return fmt.Errorf("no audio formats available for video %s", id)
	}

	audioFormats.Sort()

	stream, _, err := c.ytd.GetStreamContext(ctx, vid, &audioFormats[0])
	if err != nil {
		return err
	}

	defer stream.Close()

	title := formatName(vid.Title)
	fileName := title + ".m4a"

	if err := stream2File(stream, fileName); err != nil {
		return err
	}

	return addMetaData(ctx, vid, fileName)
}
