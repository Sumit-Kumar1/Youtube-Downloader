package handler

import (
	"testing"
	"ytdl_http/models"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTest(t *testing.T) (*Handler, *echo.Echo) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	s := NewMockServicer(ctrl)
	h := New(s)

	e.Renderer = models.NewTemplate("html")
	return h, e
}

func TestHandler_Page(t *testing.T) {
	h, e := newTest(t)
	tests := []struct {
		name    string
		wantErr error
	}{
		{"valid case", nil},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantErr, h.Page(e.AcquireContext()), "TEST[%d] Failed - %s", i, tt.name)
		})
	}
}

func TestHandler_Status(t *testing.T) {
	type fields struct {
		servicer Servicer
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				servicer: tt.fields.servicer,
			}
			if err := h.Status(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Handler.Status() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_GetInfo(t *testing.T) {
	type fields struct {
		servicer Servicer
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				servicer: tt.fields.servicer,
			}
			if err := h.GetInfo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Handler.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_Download(t *testing.T) {
	type fields struct {
		servicer Servicer
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				servicer: tt.fields.servicer,
			}
			if err := h.Download(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Handler.Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_DownloadInfo(t *testing.T) {
	type fields struct {
		servicer Servicer
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				servicer: tt.fields.servicer,
			}
			if err := h.DownloadInfo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Handler.DownloadInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
