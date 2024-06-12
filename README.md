# Iptv Proxy

[![Actions Status](https://github.com/segadora/iptv-proxy/workflows/CI/badge.svg)](https://github.com/segadora/iptv-proxy/actions?query=workflow%3ACI)

## Description

Iptv-Proxy is a project to proxyfie an m3u file

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
 - convert original m3u file with new routes pointing to the proxy

## Docker compose example with nordvpn

