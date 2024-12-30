package client

import (
	"context"
	"io"

	"github.com/kkdai/youtube/v2"
)

type ytdlr interface {
	GetPlaylist(url string) (*youtube.Playlist, error)
	GetVideo(url string) (*youtube.Video, error)
	GetStreamContext(ctx context.Context, video *youtube.Video, format *youtube.Format) (io.ReadCloser, int64, error)
	GetVideoContext(ctx context.Context, url string) (*youtube.Video, error)

	DownloadComposite(ctx context.Context, outputFile string, v *youtube.Video, quality string, mimetype, language string) error
}
