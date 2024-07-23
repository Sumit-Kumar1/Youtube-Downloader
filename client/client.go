package client

import (
	"context"
	"io"
	"os"
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

func (c Client) GetDownloadInfo(videoID string) ([]string, error) {
	vid, err := c.Ytdl.GetVideo(videoID)
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

func (c Client) DownloadVideo(id, qual string) error {
	vid, err := c.Ytdl.GetVideo(id)
	if err != nil {
		return err
	}

	title := formatName(vid.Title)
	if err := c.Ytdl.DownloadComposite(context.Background(), title+".mp4", vid, qual, "", ""); err != nil {
		return err
	}

	return nil
}

func (c Client) DownloadAudio(id string) error {
	ctx := context.Background()
	vid, err := c.Ytdl.GetVideoContext(ctx, id)
	if err != nil {
		return err
	}

	title := formatName(vid.Title)

	outFile, err := os.Create("./Downloads/" + title + ".m4a")
	if err != nil {
		return err
	}

	audioFormats := vid.Formats.Type("audio")
	audioFormats.Sort()

	stream, _, err := c.Ytdl.GetStreamContext(ctx, vid, &audioFormats[0])
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, stream)
	if err != nil {
		return err
	}

	if err := outFile.Sync(); err != nil {
		return err
	}

	return nil
}
