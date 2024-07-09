package main

import (
	"net/http"
	"os"
	"time"
	"ytdl_http/handler"
	"ytdl_http/models"
	"ytdl_http/service"

	yt "github.com/kkdai/youtube/v2"
	dlr "github.com/kkdai/youtube/v2/downloader"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// pprof.Register(e)

	t := models.NewTemplate("html")
	e.Renderer = t

	h := initServices()

	err := os.Mkdir("Downloads", 0o755)
	if err != nil && err.Error() != "mkdir Downloads: file exists" {
		e.Logger.Errorf("Error while creating folder: %s", err.Error())
		return
	}

	addMiddleWares(e)

	e.GET("/", h.Page)
	e.GET("/player", h.PagePlayer)
	e.GET("/play", h.Play)

	e.GET("/status", h.Status)
	e.GET("/info", h.DownloadInfo)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.POST("/getInfo", h.GetInfo)
	e.POST("/download", h.Download)
	e.GET("/resource/*", echo.WrapHandler(http.StripPrefix("/resource/", http.FileServer(http.Dir("./Downloads")))))

	e.Logger.Fatal(e.Start(":12344"))
}

func initServices() *handler.Handler {
	c := &yt.Client{MaxRoutines: 8, ChunkSize: yt.Size1Mb}
	d := &dlr.Downloader{OutputDir: "Downloads"}
	s := service.New(c, d)

	return handler.New(s)
}

func addMiddleWares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("ytdl_http"))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
	// 	rate.Limit(160),
	// )))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timed-out due to no activity for 3 mins",
		Timeout:      180 * time.Second,
	}))
}
