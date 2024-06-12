FROM golang:1.22-alpine as build

WORKDIR /app

COPY *.go ./
COPY go.mod go.sum ./

RUN ls -lA

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o iptv-proxy *.go

FROM alpine:3 as executeable

COPY --from=build  /app/iptv-proxy /

ENTRYPOINT ["/iptv-proxy"]
