package main

import (
	"log"
	stdhttp "net/http"
	"path/filepath"
	"time"

	routerpkg "http_server/internal/adapter/http"
	handlerpkg "http_server/internal/adapter/http/handler"
	ext "http_server/internal/external/http"
	"http_server/internal/service"
)

func main() {
    httpClient := &stdhttp.Client{Timeout: 10 * time.Second}

    weatherClient := ext.NewWeatherClient(httpClient)
    geoClient := ext.NewGeolocationClient(httpClient)

    svc := service.NewWeatherService(weatherClient, geoClient)

    h := handlerpkg.NewHandler(svc)

    // Path to openapi.yaml (relative to cmd/api working directory).
    openapiPath := filepath.Join("..", "..", "api", "openapi.yaml")
    r := routerpkg.NewRouter(h, openapiPath)

    addr := ":8080"
    log.Printf("starting server on %s", addr)
    if err := stdhttp.ListenAndServe(addr, r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
