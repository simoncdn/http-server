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
	resBody := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, hits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}
