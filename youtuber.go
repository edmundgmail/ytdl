package main

import (
	"context"
	"strconv"
	"net/http"
	"os"
	"github.com/Andreychik32/ytdl"
	"github.com/labstack/echo/v4"
	"encoding/base64"
)

type youtuber struct {
	e *echo.Echo
	cxt context.Context
	client *ytdl.Client
	tracer Tracer
}

func (r *youtuber) DownloadVideo(c echo.Context) (err error) {
	encoded := c.Param("url")
	id, _ := strconv.Atoi(c.Param("id"))
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	url := string(decoded);

	videoInfo := r.getVideoInfo(url)
	format := videoInfo.Formats[id]
	filename := videoInfo.Title+"."+format.Extension
	
	res := c.Response()
	header := res.Header()
	header.Set(echo.HeaderContentType, echo.MIMEOctetStream)
	header.Set(echo.HeaderContentDisposition, "attachment; filename="+filename)
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	res.WriteHeader(http.StatusOK)

	err1 := r.client.Download(r.cxt, videoInfo, format, res)
	if err1 != nil {
		return
	}

	res.Flush()
	return
}

func (r *youtuber) GetVideoInfo(c echo.Context) error {
	encoded := c.Param("url")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err)
	}

	url := string(decoded)
	return c.JSON(http.StatusOK, *(r.getVideoInfo(url)))
}

func (r *youtuber) getVideoInfo(url string) *ytdl.VideoInfo{
	videoInfo, err := r.client.GetVideoInfo(r.cxt, url)
	if err != nil {
		panic(err)
	}
	return videoInfo
}

func newYoutuber(e *echo.Echo) *youtuber {
	return &youtuber{
		e: e,
		cxt: context.Background(),
		client: ytdl.DefaultClient,
		tracer: New(os.Stdout),
	}
}
