package server

import (
	"fmt"
	"net/http"

	"github.com/simoncdn/http-server/internal/config"
	"github.com/simoncdn/http-server/internal/handlers"
	"github.com/simoncdn/http-server/internal/middleware"
)

type Server struct {
	config *config.Config
	mux    *http.ServeMux
	server *http.Server
}

func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	server := &Server{
		config: cfg,
		mux:    mux,
		server: httpServer,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	metricsMiddleware := middleware.Metrics(s.config)

	fileServer := http.FileServer(http.Dir(s.config.StaticDir))
	fileServerHandler := http.StripPrefix("/app", fileServer)

	metricsHandler := handlers.NewMetricsHandler(s.config)

	s.mux.Handle("/app/", metricsMiddleware(fileServerHandler))
	s.mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	s.mux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirp)
	s.mux.HandleFunc("GET /admin/metrics", metricsHandler.GetMetrics)
	s.mux.HandleFunc("POST /admin/reset", metricsHandler.ResetMetrics)
}

func (s *Server) Start() error {
	fmt.Printf("serving on port: %s\n", s.config.Port)
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error on serving: %w", err)
	}

	return nil
}
