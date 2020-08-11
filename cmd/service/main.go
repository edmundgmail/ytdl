package main

import (

	"flag"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())	

	var addr = flag.String("addr", ":8000", "The addr of the application.")
	flag.Parse() // parse the flags

	http.Handle("/", &templateHandler{filename: "dl.html"})

	r := newYoutuber(e)
	e.GET("/info/:url", r.GetVideoInfo)
	e.GET("/download/:url/:id", r.DownloadVideo)
	//go r.run()
	e.Logger.Fatal(e.Start(*addr))
}
