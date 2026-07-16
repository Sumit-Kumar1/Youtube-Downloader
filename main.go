package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ytdl_http/internal/client"
	"ytdl_http/internal/handler"
	"ytdl_http/internal/models"
	"ytdl_http/internal/service"

	echoprometheus "github.com/labstack/echo-prometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	t := models.NewTemplate("html")

	e.Renderer = t

	if _, err := os.Stat(models.DirPath); os.IsNotExist(err) {
		if err := os.Mkdir(models.DirPath, 0o750); err != nil {
			e.Logger.Error("directory creation error", slog.String("error", err.Error()), slog.String("dir-path", models.DirPath))

			return
		}
	}

	h := setupDeps()

	addMiddlewares(e)

	e.GET("/", h.Page)
	e.GET("/player", h.PagePlayer)
	e.GET("/play", h.Play)

	e.GET("/info", h.DownloadInfo)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.POST("/getInfo", h.GetInfo)
	e.POST("/download", h.Download)

	e.GET("/resource/*", echo.WrapHandler(http.StripPrefix("/resource/",
		http.FileServer(http.Dir(models.DirPath)))))
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("assets")))))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":9001",
		GracefulTimeout: 5 * time.Second,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", slog.String("error", err.Error()))
	}
}

func setupDeps() *handler.Handler {
	ytCl := client.New()
	s := service.New(ytCl)
	return handler.New(s)
}

func addMiddlewares(app *echo.Echo) {
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(echoprometheus.NewMiddleware(models.AppName))
	app.Use(middleware.Recover())
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(150.0)))
}
