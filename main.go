package main

import (
	"net/http"
	"os"
	"time"

	"ytdl_http/internal/client"
	"ytdl_http/internal/handler"
	"ytdl_http/internal/models"
	"ytdl_http/internal/service"

	"golang.org/x/time/rate"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	t := models.NewTemplate("html")

	e.Renderer = t

	if _, err := os.Stat(models.DirPath); os.IsNotExist(err) {
		if err := os.Mkdir(models.DirPath, 0755); err != nil {
			e.Logger.Errorf("error while creating the directory: %s", err.Error())

			return
		}
	}

	h := initServices()

	addMiddleWares(e)

	e.GET("/", h.Page)
	e.GET("/player", h.PagePlayer)
	e.GET("/play", h.Play)

	e.GET("/info", h.DownloadInfo)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.POST("/getInfo", h.GetInfo)
	e.POST("/download", h.Download)
	e.GET("/resource/*", echo.WrapHandler(http.StripPrefix("/resource/", http.FileServer(http.Dir("./Downloads")))))

	e.Logger.Fatal(e.Start(":12344"))
}

func initServices() *handler.Handler {
	ytCl := client.New()
	s := service.New(ytCl)

	return handler.New(s)
}

func addMiddleWares(app *echo.Echo) {
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Logger())
	app.Use(echoprometheus.NewMiddleware("ytdl_http"))
	app.Use(middleware.Recover())
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(160),
	)))
	app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timed-out due to no activity for 3 mins",
		Timeout:      180 * time.Second,
	}))
}
