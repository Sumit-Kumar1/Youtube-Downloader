package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
	"ytdl_http/handler"
	"ytdl_http/models"
	"ytdl_http/service"

	"golang.org/x/time/rate"

	yt "github.com/kkdai/youtube/v2"
	dlr "github.com/kkdai/youtube/v2/downloader"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//profiling
	cpuProfFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("not able to create a profiling file")
		return
	}

	err = pprof.StartCPUProfile(cpuProfFile)
	if err != nil {
		fmt.Println("not able to start profiling")
		return
	}

	defer func() {
		pprof.StopCPUProfile()
		cpuProfFile.Close()
	}()

	e := echo.New()
	t := models.NewTemplate("html")
	e.Renderer = t

	h := initServices()

	err = os.Mkdir("Downloads", 0o755)
	if err != nil && err.Error() != "mkdir Downloads: file exists" {
		e.Logger.Errorf("Error while creating folder: %s", err.Error())
		return
	}

	addMiddleWares(e)

	// routes
	e.GET("/", h.Page)
	e.GET("/status", h.Status)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.POST("/getInfo", h.GetInfo)
	e.POST("/download", h.Download)
	e.GET("/info", h.DownloadInfo)

	e.Logger.Fatal(e.Start(":12344"))
}

func initServices() *handler.Handler {
	c := &yt.Client{MaxRoutines: 8, ChunkSize: yt.Size10Mb}
	d := &dlr.Downloader{OutputDir: "Downloads"}
	s := service.New(c, d)

	return handler.New(s)
}

func addMiddleWares(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoprometheus.NewMiddleware("youtube_downloader"))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(10),
	)))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timed-out due to no activity for 3 mins",
		Timeout:      180 * time.Second,
	}))
}
