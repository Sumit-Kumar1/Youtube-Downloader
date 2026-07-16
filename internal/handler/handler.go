package handler

import (
	"log/slog"
	"net/http"

	"ytdl_http/internal/models"

	"github.com/labstack/echo/v5"
)

// Handler type contains service type for dependency injection.
type Handler struct {
	service Servicer
}

// New to init the handler.
func New(s Servicer) *Handler {
	return &Handler{
		service: s,
	}
}

// Page render the root index.html page.
func (h *Handler) Page(c *echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

// PagePlayer render the player page.
func (h *Handler) PagePlayer(c *echo.Context) error {
	data, err := getDownloadedFiles()
	if err != nil {
		c.Logger().Error("player page: get download files error", slog.String("error", err.Error()))
		return c.Render(http.StatusOK, "error", map[string]string{"error": err.Error()})
	}

	return c.Render(http.StatusOK, "Player", map[string]any{models.Data: data})
}

func (h *Handler) Play(c *echo.Context) error {
	var p models.Player

	title := c.QueryParam(models.Title)

	p.FillByName(title)

	switch p.Type {
	case models.AudioType:
		return c.Render(http.StatusOK, "audio", map[string]any{models.Data: p})
	case models.VideoType:
		return c.Render(http.StatusOK, "video", map[string]any{models.Data: p})
	default:
		return c.Render(http.StatusOK, "error", map[string]string{"error": "unsupported media type"})
	}
}

// GetInfo retrieves the information about the url provided.
func (h *Handler) GetInfo(c *echo.Context) error {
	url := c.FormValue("URL")

	data, err := h.service.GetInfo(c.Request().Context(), url)
	if err != nil {
		return c.Render(http.StatusOK, "error", map[string]string{"error": err.Error()})
	}

	d := map[string][]models.Video{models.Data: data}

	return c.Render(http.StatusOK, "info", d)
}

// Download start download process based on quality and videoID.
func (h *Handler) Download(c *echo.Context) error {
	qual := c.FormValue("quality")
	id := c.FormValue("id")
	audioOnly := c.FormValue("audioOnly") == "true"

	if err := h.service.Download(c.Request().Context(), id, qual, audioOnly); err != nil {
		c.Logger().Error("download error", slog.String("error", err.Error()))
		return c.Render(http.StatusOK, "error", map[string]string{"error": "download failed, please try again"})
	}

	return c.Render(http.StatusOK, "status", map[string]any{"ID": id})
}

func (h *Handler) DownloadInfo(c *echo.Context) error {
	videoID := c.QueryParam("id")

	qualities, err := h.service.DownloadInfo(c.Request().Context(), videoID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "downloadInfo", map[string]any{
		"ID":         videoID,
		"VidQuality": qualities,
	})
}
