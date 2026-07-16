package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"ytdl_http/internal/models"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

var errService = errors.New("service error")

type stubServicer struct {
	getInfoFn      func(ctx context.Context, url string) ([]models.Video, error)
	downloadInfoFn func(ctx context.Context, videoID string) ([]string, error)
	downloadFn     func(ctx context.Context, id, qual string, audioOnly bool) error
}

func (s *stubServicer) GetInfo(ctx context.Context, url string) ([]models.Video, error) {
	if s.getInfoFn != nil {
		return s.getInfoFn(ctx, url)
	}

	return nil, nil
}

func (s *stubServicer) DownloadInfo(ctx context.Context, videoID string) ([]string, error) {
	if s.downloadInfoFn != nil {
		return s.downloadInfoFn(ctx, videoID)
	}

	return nil, nil
}

func (s *stubServicer) Download(ctx context.Context, id, qual string, audioOnly bool) error {
	if s.downloadFn != nil {
		return s.downloadFn(ctx, id, qual, audioOnly)
	}

	return nil
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Renderer = models.NewTemplate("../../html")

	return e
}

func TestHandler_Page(t *testing.T) {
	e := setupEcho()
	h := New(&stubServicer{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Page(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_GetInfo(t *testing.T) {
	e := setupEcho()

	tests := []struct {
		name       string
		url        string
		stub       *stubServicer
		wantStatus int
	}{
		{
			name: "success",
			url:  "https://www.youtube.com/watch?v=abc123",
			stub: &stubServicer{
				getInfoFn: func(_ context.Context, _ string) ([]models.Video, error) {
					return []models.Video{{ID: "abc123", Title: "Test"}}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error renders error template",
			url:  "invalid",
			stub: &stubServicer{
				getInfoFn: func(_ context.Context, _ string) ([]models.Video, error) {
					return nil, errService
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.stub)

			form := url.Values{}
			form.Set("URL", tt.url)

			req := httptest.NewRequest(http.MethodPost, "/getInfo", strings.NewReader(form.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GetInfo(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandler_Download(t *testing.T) {
	e := setupEcho()

	tests := []struct {
		name         string
		formValues   map[string]string
		stub         *stubServicer
		wantStatus   int
		wantContains string
	}{
		{
			name: "successful download",
			formValues: map[string]string{
				"id": "abc123", "quality": "720p", "audioOnly": "",
			},
			stub: &stubServicer{
				downloadFn: func(_ context.Context, id, qual string, audioOnly bool) error {
					assert.Equal(t, "abc123", id)
					assert.Equal(t, "720p", qual)
					assert.False(t, audioOnly)

					return nil
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: "Download Completed",
		},
		{
			name: "audio only download",
			formValues: map[string]string{
				"id": "abc123", "quality": "", "audioOnly": "true",
			},
			stub: &stubServicer{
				downloadFn: func(_ context.Context, _ string, _ string, audioOnly bool) error {
					assert.True(t, audioOnly)

					return nil
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: "Download Completed",
		},
		{
			name: "download error shows user-friendly message",
			formValues: map[string]string{
				"id": "abc123", "quality": "720p", "audioOnly": "",
			},
			stub: &stubServicer{
				downloadFn: func(_ context.Context, _, _ string, _ bool) error {
					return errService
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: "download failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.stub)

			form := url.Values{}
			for k, v := range tt.formValues {
				form.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodPost, "/download", strings.NewReader(form.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.Download(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestHandler_Play_UnsupportedType(t *testing.T) {
	e := setupEcho()
	h := New(&stubServicer{})

	req := httptest.NewRequest(http.MethodGet, "/play?title=readme.txt", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Play(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "unsupported media type")
}

func TestHandler_DownloadInfo(t *testing.T) {
	e := setupEcho()

	tests := []struct {
		name       string
		videoID    string
		stub       *stubServicer
		wantStatus int
	}{
		{
			name:    "success",
			videoID: "abc123",
			stub: &stubServicer{
				downloadInfoFn: func(_ context.Context, _ string) ([]string, error) {
					return []string{"720p", "480p", "360p"}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "error",
			videoID: "invalid",
			stub: &stubServicer{
				downloadInfoFn: func(_ context.Context, _ string) ([]string, error) {
					return nil, errService
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.stub)

			req := httptest.NewRequest(http.MethodGet, "/info?id="+tt.videoID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.DownloadInfo(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
