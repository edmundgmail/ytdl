package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())	

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	r := newYoutuber(e)
	e.File("/", "public/index.html")
	e.GET("/info/:url", r.GetVideoInfo)
	e.GET("/download/:url/:id", r.DownloadVideo)
	e.Logger.Fatal(e.Start(*addr))
}
