FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go config.go ./
COPY static ./static
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/neonvr .
FROM alpine:3.20
RUN apk add --no-cache ffmpeg ca-certificates tzdata
COPY --from=build /out/neonvr /usr/local/bin/neonvr
EXPOSE 5500
ENTRYPOINT ["/usr/local/bin/neonvr"]
CMD ["-config", "/config/neonvr.json"]
