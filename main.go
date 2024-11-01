package main

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/jamesnetherton/m3u"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	Playlist  string
	EPG       string
	ServerUrl string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("unable to load .env file: %s", err)
	}

	config := &Config{
		Playlist:  os.Getenv("IPTV_PLAYLIST"),
		EPG:       os.Getenv("IPTV_EPG"),
		ServerUrl: os.Getenv("IPTV_SERVER_URL"),
	}

	log.Printf("playlist url: %s", config.Playlist)
	log.Printf("epg url: %s", config.EPG)
	log.Printf("server url: %s", config.ServerUrl)
	if config.ServerUrl == "" {
		log.Printf("warning: host is dynamically defined on request")
	}

	r := gin.New()
	r.Use(
		cors.Default(),
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/health"}}),
		gin.Recovery(),
	)
	r.GET("/health", healthHandler)
	r.GET("/get/epg", config.epgHandler)
	r.GET("/get/m3u", config.playlistHandler)
	r.POST("/get/epg", config.epgHandler)
	r.POST("/get/m3u", config.playlistHandler)
	r.GET("/stream", logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			var streamRequest StreamRequest
			if c.ShouldBind(&streamRequest) != nil {
				return l.With().
					Str("channel", "unknown").
					Logger()
			}

			return l.With().
				Str("channel", streamRequest.ChannelName).
				Logger()
		}),
	), streamHandler)

	if err := r.Run(":1323"); err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

type StreamRequest struct {
	RemoteUrl   string `form:"remoteUrl"`
	ChannelName string `form:"channelName"`
}

func streamHandler(c *gin.Context) {
	var streamRequest StreamRequest
	if err := c.ShouldBind(&streamRequest); err != nil {
		log.Printf("bad request: %s", err)

		c.String(http.StatusBadRequest, "bad request")

		return
	}

	req, err := http.NewRequest("GET", streamRequest.RemoteUrl, nil)
	if err != nil {
		log.Printf("error creating HTTP request: %s", err.Error())

		_ = c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	mergeHttpHeader(req.Header, c.Request.Header)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error creating HTTP request: %s", err.Error())

		_ = c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("error closing response body: %s", err.Error())
		}
	}(resp.Body)

	mergeHttpHeader(c.Writer.Header(), resp.Header)

	c.Status(resp.StatusCode)
	c.Stream(func(w io.Writer) bool {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return false
		}
		return false
	})
}

func (config *Config) epgHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, config.EPG)
}

func (config *Config) playlistHandler(c *gin.Context) {
	serverUrl := config.ServerUrl
	if serverUrl == "" {
		serverUrl = fmt.Sprintf("%s://%s", "http", c.Request.Host)
	}

	p, err := m3u.Parse(config.Playlist)
	if err != nil {
		log.Printf("unable to parse playlist: %s", err)

		c.String(http.StatusInternalServerError, "unable to parse playlist")

		return
	}

	if len(p.Tracks) == 0 {
		log.Printf("unable to parse playlist: playlist is empty")

		c.String(http.StatusInternalServerError, "playlist is empty")

		return
	}

	c.Header("Content-Disposition", "attachment; filename=playlist.m3u")
	c.Header("Content-Type", "application/octet-stream")

	if _, err := c.Writer.WriteString("#EXTM3U\n"); err != nil {
		return
	}

	for _, track := range p.Tracks {
		var buffer bytes.Buffer
		buffer.WriteString("#EXTINF:")
		buffer.WriteString(fmt.Sprintf("%d ", track.Length))
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				buffer.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
				continue
			}
			buffer.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
		}

		line := fmt.Sprintf(
			"%s, %s\n%s\n",
			buffer.String(),
			track.Name,
			serverUrl+"/stream?remoteUrl="+url.QueryEscape(track.URI)+"&channelName="+url.QueryEscape(track.Name),
		)

		if _, err := c.Writer.WriteString(line); err != nil {
			return
		}
	}
}

type values []string

func (vs values) contains(s string) bool {
	for _, v := range vs {
		if v == s {
			return true
		}
	}

	return false
}

func mergeHttpHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			if values(dst.Values(k)).contains(v) {
				continue
			}
			dst.Add(k, v)
		}
	}
}
