package main

import (
	"context"
	"strconv"
	"encoding/json"
	"net/http"
	"os"
	"io"
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

type UrlRequest struct {
	Link string `json:"link"`	
	Id int `json:"id"`
}



type VideoFormat struct {
	Itag int `json: "itag"`
	Extension string `json: "Extension"`
	AudioEncoding string `json: "AudioEncoding"`
	AudioBitrate string `json: "AudioBitrate"`
	VideoEncoding string `json: "VideoEncoding"`
	Resolution string `json: "Resolution"`
}

type VideoInformation struct {
	Title string `json:"Title"`
	Formats []VideoFormat `json:"Formats"`
}

func (r *youtuber) downloadVideo(url string, index int, out io.Writer) (string, string) {
	videoInfo := r.getVideoInfo(url)
	filename := videoInfo.ID+"."+videoInfo.Formats[index].Extension
	title := videoInfo.Title +"."+videoInfo.Formats[index].Extension
	err := r.client.Download(r.cxt, videoInfo, videoInfo.Formats[index], out)
	if err != nil {
		panic(err)
	}

	return title, filename
}

func (r *youtuber) DownloadVideo(c echo.Context) (err error) {
	encoded := c.Param("url")
	id, _ := strconv.Atoi(c.Param("id"))
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return
	}
	res := c.Response()
	header := res.Header()
	header.Set(echo.HeaderContentType, echo.MIMEOctetStream)
	header.Set(echo.HeaderContentDisposition, "attachment; filename="+"abc.mp4")
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	res.WriteHeader(http.StatusOK)

	r.downloadVideo(string(decoded), id, res)
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

func (r *youtuber) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if(req.Method=="POST"){
		decoder := json.NewDecoder(req.Body)
		var t UrlRequest
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		r.getVideoInfo(t.Link)
	}
		
}

func newYoutuber(e *echo.Echo) *youtuber {
	return &youtuber{
		e: e,
		cxt: context.Background(),
		client: ytdl.DefaultClient,
		tracer: New(os.Stdout),
	}
}
