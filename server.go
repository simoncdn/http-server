package main

import (
	"fmt"
	"net/http"
)

const (
	port         = "8080"
	filepathRoot = "."
)

func server(cfg *apiConfig) error {
	fileServer := http.FileServer(http.Dir(filepathRoot))
	fileServerHandler := http.StripPrefix("/app", fileServer)

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("/reset", cfg.resetMetricsHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("serving on port: %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error on serving: %w", err)
	}

	return nil
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	resBody := fmt.Sprintf("Hits: %d", hits)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	resBody := "Hits reset to 0"

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))
}
