package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheusExporter() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))

	port := os.Getenv("PROMETHEUS_PORT")
	if port == "" {
		port = "3002"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	go func() {
		log.Info("Starting Prometheus exporter ", "port ", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Prometheus exporter failed", "err", err)
		}
	}()
}
