package handler

import (
	"net/http"

	"ytdl_http/internal/models"

	"github.com/labstack/echo/v4"
)

// Handler type contains service type for dependency injection.
type Handler struct {
	player  map[string]bool
	service Servicer
}

// New to init the handler.
func New(s Servicer) *Handler {
	return &Handler{
		player:  make(map[string]bool, 0),
		service: s,
	}
}

// Page render the root index.html page.
//
//nolint:revive // to make it consistent at main
func (h *Handler) Page(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

// Page render the root index.html page.
//
//nolint:revive // to make it consistent at main
func (h *Handler) PagePlayer(c echo.Context) error {
	data, err := getDownloadedFiles()
	if err != nil {
		c.Logger().Errorf(err.Error())
		return c.Render(http.StatusOK, "error", map[string]string{"error": ""})
	}

	return c.Render(http.StatusOK, "Player", map[string]any{"Data": data})
}

//nolint:revive // to make it consistent at main
func (h *Handler) Play(c echo.Context) error {
	var p models.Player

	title := c.QueryParam("title")

	p.FillByName(title)

	switch p.Type {
	case "audio/mpeg":
		return c.Render(http.StatusOK, "audio", map[string]any{"Data": p})
	case "video/mp4":
		return c.Render(http.StatusOK, "video", map[string]any{"Data": p})
	default:
	}

	return nil
}

// GetInfo retrives the information about the url provided.
func (h *Handler) GetInfo(c echo.Context) error {
	url := c.FormValue("URL")

	data, err := h.service.GetInfo(url)
	if err != nil {
		return c.Render(http.StatusOK, "error", map[string]string{"error": err.Error()})
	}

	d := map[string][]models.Video{
		"Data": data,
	}

	return c.Render(http.StatusOK, "info", d)
}

// Download start download process based on quality and videoID
func (h *Handler) Download(c echo.Context) error {
	qual := c.FormValue("quality")
	id := c.FormValue("id")
	audioOnly := c.FormValue("audioOnly")

	if err := h.service.Download(id, qual, audioOnly); err != nil {
		c.Logger().Errorf("Download err: %s", err.Error())

		return err
	}

	return c.Render(http.StatusOK, "status", map[string]any{"ID": id})
}

func (h *Handler) DownloadInfo(c echo.Context) error {
	videoID := c.QueryParam("id")

	qualities, err := h.service.DownloadInfo(videoID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "downloadInfo", map[string]any{
		"ID":         videoID,
		"VidQuality": qualities,
	})
}
