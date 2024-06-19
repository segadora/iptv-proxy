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
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("unable to load .env file: %s", err)
	}

	config := &Config{
		M3UUrl: os.Getenv("M3U_URL"),
		EPGUrl: os.Getenv("EPG_URL"),
	}

	log.Printf("M3U URL: %s", config.M3UUrl)
	log.Printf("EPG URL: %s", config.EPGUrl)

	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if e := server.Serve(); e != nil {
		log.Fatal(e)
	}
}
