package handlers

import (
	"fmt"
	"net/http"

	"github.com/simoncdn/http-server/internal/config"
)

type MetricsHandler struct {
	config *config.Config
}

func NewMetricsHandler(cfg *config.Config) *MetricsHandler {
	metricsHandler := &MetricsHandler{
		config: cfg,
	}

	return metricsHandler
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	hits := h.config.GetHits()
	resBody := fmt.Sprintf("Hits: %d", hits)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}

func (h *MetricsHandler) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	h.config.ResetHits()
	resBody := "Hits reset to 0"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}
