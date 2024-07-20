package main

import (
	"net/http"
	"time"
	"ytdl_http/client"
	"ytdl_http/handler"
	"ytdl_http/models"
	"ytdl_http/service"

	dlr "github.com/kkdai/youtube/v2/downloader"
	"golang.org/x/time/rate"

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
	d := &dlr.Downloader{OutputDir: "/tmp/Downloads"}
	ytCl := client.New(d)
	s := service.New(ytCl)

	return handler.New(s)
}

func addMiddleWares(app *echo.Echo) {
	app.Use(middleware.Logger())
	app.Use(echoprometheus.NewMiddleware("ytdl_http"))
	app.Use(middleware.Recover())
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(160),
	)))

	app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timed-out due to no activity for 3 mins",
		Timeout:      180 * time.Second,
	}))
}
