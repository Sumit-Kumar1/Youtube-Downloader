package handler

import (
	"downloader/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler type contains service type for dependency injection.
type Handler struct {
	servicer Servicer
}

// New to init the handler.
func New(s Servicer) *Handler {
	return &Handler{servicer: s}
}

// Page render the root index.html page.
func (h *Handler) Page(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

// Status returns the download status of a video
func (h *Handler) Status(c echo.Context) error {
	vID := c.QueryParam("vid")
	status := h.servicer.GetStatus(vID)

	return c.String(http.StatusOK, status)
}

// GetInfo retrives the information about the url provided.
func (h *Handler) GetInfo(c echo.Context) error {
	url := c.FormValue("URL")
	if url == "" {
		return c.Render(http.StatusOK, "error", map[string]string{"error": "url is empty"})
	}

	data, err := h.servicer.GetInfo(url)
	if err != nil {
		return c.Render(http.StatusOK, "error", map[string]string{"error": err.Error()})
	}

	resp := map[string][]models.VideoData{
		"Data": data,
	}

	if len(data) > 1 {
		return c.Render(http.StatusOK, "playlistInfo", resp)
	}

	return c.Render(http.StatusOK, "videoInfo", resp)
}

// Download start download process based on quality and videoID
func (h *Handler) Download(c echo.Context) error {
	qual := c.FormValue("quality")
	id := c.FormValue("id")
	audioOnly := c.FormValue("audioOnly")

	fmt.Println("Audio only: ", audioOnly)

	switch audioOnly {
	case "":
		if err := h.servicer.DownloadVideo(c.Request().Context(), id, qual); err != nil {
			return c.String(http.StatusBadRequest, "Error h.DownloadVideo: "+err.Error())
		}
	case "true":
		if err := h.servicer.DownloadAudio(c.Request().Context(), id); err != nil {
			return c.String(http.StatusBadRequest, "Error h.DownloadAudio: "+err.Error())
		}
	}

	css := "text-white bg-blue-700 hover:bg-blue-800 focus:outline-none focus:ring-4 focus:ring-blue-300 font-medium rounded-md text-sm px-5 py-2.5 text-center me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
	msg := fmt.Sprintf(`<p id='status' class="text-green-500">Download Started!!</p> <button type='button' hx-get='/status?vid=%s' hx-target="#status" class='%s'>Get Status</button>`, id, css)

	return c.String(http.StatusOK, msg)
}

func (h *Handler) DownloadInfo(c echo.Context) error {
	vidID := c.QueryParam("id")
	data, err := h.servicer.DownloadInfo(vidID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Render(http.StatusOK, "downloadInfo", map[string]interface{}{
		"VideoID":    vidID,
		"VidQuality": data,
	})
}
