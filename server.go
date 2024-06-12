/*
 * Iptv-Proxy is a project to proxyfie an m3u file and to proxyfie an Xtream iptv service (client API).
 * Copyright (C) 2020  Pierre-Emmanuel Jacquier
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jamesnetherton/m3u"
	"log"
	"net/url"
	"os"
	"path"
)

// Server represent the server configuration
type Server struct {
	*Config

	proxyPlaylistFile *os.File

	// M3U service part
	playlist *m3u.Playlist
	// this variable is set only for m3u proxy endpoints
	track *m3u.Track
}

// NewServer initialize a new server configuration
func NewServer(config *Config) (*Server, error) {
	var p m3u.Playlist
	var err error
	p, err = m3u.Parse(config.M3UUrl)
	if err != nil {
		return nil, err
	}

	return &Server{
		config,
		nil,
		&p,
		nil,
	}, nil
}

// Serve the iptv-proxy api
func (c *Server) Serve() error {
	if err := c.playlistInitialization(); err != nil {
		return err
	}

	router := gin.Default()
	router.Use(cors.Default())
	group := router.Group("/")
	c.routes(group)

	return router.Run(fmt.Sprintf(":%d", 1323))
}

func (c *Server) playlistInitialization() error {
	if len(c.playlist.Tracks) == 0 {
		return nil
	}

	f, err := os.CreateTemp("", "playlist")
	if err != nil {
		return err
	}
	defer f.Close()

	c.proxyPlaylistFile = f

	return c.marshallInto(f)
}

// MarshallInto a *bufio.Writer a Playlist.
func (c *Server) marshallInto(into *os.File) error {
	filteredTrack := make([]m3u.Track, 0, len(c.playlist.Tracks))

	ret := 0
	into.WriteString("#EXTM3U\n") // nolint: errcheck
	for i, track := range c.playlist.Tracks {
		var buffer bytes.Buffer

		buffer.WriteString("#EXTINF:")                       // nolint: errcheck
		buffer.WriteString(fmt.Sprintf("%d ", track.Length)) // nolint: errcheck
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				buffer.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
				continue
			}
			buffer.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
		}

		uri, err := c.replaceURL(track.URI, i-ret)
		if err != nil {
			ret++
			log.Printf("ERROR: track: %s: %s", track.Name, err)
			continue
		}

		into.WriteString(fmt.Sprintf("%s, %s\n%s\n", buffer.String(), track.Name, uri)) // nolint: errcheck

		filteredTrack = append(filteredTrack, track)
	}
	c.playlist.Tracks = filteredTrack

	return into.Sync()
}

// ReplaceURL replace original playlist url by proxy url
func (c *Server) replaceURL(uri string, trackIndex int) (string, error) {
	oriURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	uriPath := oriURL.EscapedPath()
	uriPath = path.Join(fmt.Sprintf("%d", trackIndex), path.Base(uriPath))

	newURI := fmt.Sprintf(
		"/%s",
		uriPath,
	)

	newURL, err := url.Parse(newURI)
	if err != nil {
		return "", err
	}

	return "[URL]" + newURL.String(), nil
}
