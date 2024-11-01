# Iptv Proxy

[![Actions Status](https://github.com/segadora/iptv-proxy/workflows/CI/badge.svg)](https://github.com/segadora/iptv-proxy/actions?query=workflow%3ACI)

## Description

Iptv-Proxy is acting as a reverse proxy and serving the stream though go to allow use of a proxy in docker compose.

### M3u Example

Original iptv m3u file

```m3u
#EXTM3U
#EXTINF:-1 tvg-ID="examplechanel1.com" tvg-name="chanel1" tvg-logo="http://ch.xyz/logo1.png" group-title="USA HD",CHANEL1-HD
http://iptvexample.net:1234/12/test/1
#EXTINF:-1 tvg-ID="examplechanel2.com" tvg-name="chanel2" tvg-logo="http://ch.xyz/logo2.png" group-title="USA HD",CHANEL2-HD
http://iptvexample.net:1234/13/test/2
#EXTINF:-1 tvg-ID="examplechanel3.com" tvg-name="chanel3" tvg-logo="http://ch.xyz/logo3.png" group-title="USA HD",CHANEL3-HD
http://iptvexample.net:1234/14/test/3
#EXTINF:-1 tvg-ID="examplechanel4.com" tvg-name="chanel4" tvg-logo="http://ch.xyz/logo4.png" group-title="USA HD",CHANEL4-HD
http://iptvexample.net:1234/15/test/4
```

What M3U proxy IPTV do
 - convert chanels url to new endpoints
 - convert original m3u file with new routes pointing to internal routes

## Docker compose example with nordvpn

Uses docker container from [edgd1er/nordvpn-proxy](https://github.com/edgd1er/nordvpn-proxy).

The following urls will be available for you.

M3U: `http://127.0.0.1:1323/get/m3u`

EPG: `http://127.0.0.1:1323/get/epg`

Health endpoint: `http://127.0.0.1:1323/health`

```yaml
services:
  proxy:
    image: edgd1er/nordvpn-proxy:latest
    restart: unless-stopped
    container_name: proxy
    ports:
      - "1081:1080"
      - "8888:8888/tcp"
    sysctls:
      - net.ipv4.conf.all.rp_filter=2
    cap_add:
      - NET_ADMIN
    environment:
      - TZ=Europe/London
      - DNS=1.1.1.1@853#cloudflare-dns.com 1.0.0.1@853#cloudflare-dns.com
      - NORDVPN_COUNTRY=Germany
      - NORDVPN_PROTOCOL=udp
      - NORDVPN_CATEGORY=p2p
      - EXIT_WHEN_IP_NOTASEXPECTED=1
      - WRITE_OVPN_STATUS=0
      - LOCAL_NETWORK=192.168.1.0/24
      - TINYPORT=8888
      - TINY_LOGLEVEL=Error
      - DANTE_LOGLEVEL="error connect"
      - DANTE_ERRORLOG=/dev/stdout
      - CRON_LOGLEVEL=9
      - DEBUG=0
    secrets:
      - NORDVPN_CREDS
    volumes:
      - ./nordvpn/config/:/config/

  iptv-proxy:
    image: ghcr.io/segadora/iptv-proxy:latest
    container_name: "iptv-proxy"
    restart: on-failure
    environment:
      HTTP_PROXY: "http://proxy:8888"
      HTTPS_PROXY: "http://proxy:8888"
      NO_PROXY: "127.0.0.0/8"
      IPTV_PLAYLIST: https://xeev.net/get/m3u/xxxxxxxxxxxxxxxxxxxxx
      IPTV_EPG: https://xeev.net/get/epg/xxxxxxxxxxxxxxxxxxxxx
    depends_on:
      proxy:
        condition: service_healthy
```
