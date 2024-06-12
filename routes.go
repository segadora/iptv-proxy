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
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
)

func (c *Server) routes(r *gin.RouterGroup) {
	r = r.Group("")

	c.m3uRoutes(r)
}

func (c *Server) m3uRoutes(r *gin.RouterGroup) {
	r.GET("/playlist.m3u", c.getM3U)
	r.POST("/playlist.m3u", c.getM3U)
	r.GET("/epg", c.getEPG)

	for i, track := range c.playlist.Tracks {
		trackConfig := &Server{
			track: &c.playlist.Tracks[i],
		}

		r.GET(fmt.Sprintf("/%d/%s", i, path.Base(track.URI)), trackConfig.reverseProxy)
	}
}
