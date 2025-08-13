package handlers

import (
	"net/http"

	"github.com/simoncdn/http-server/internal/config"
)

type ResetHandler struct {
	config *config.Config
}

func NewResetHandler(cfg *config.Config) *ResetHandler {
	resetHandler := &ResetHandler{
		config: cfg,
	}

	return resetHandler
}

func (h *ResetHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if h.config.Plateform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	h.config.ResetHits()
	resBody := "Hits reset to 0 | Users reset"

	h.config.ResetUsers(r.Context())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}
