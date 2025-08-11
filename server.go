package main

import (
	"fmt"
	"net/http"
)

const (
	port         = "8080"
	filepathRoot = "."
)

func server() error {
	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", healthzHandler)

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

func healthzHandler (res http.ResponseWriter, _ *http.Request) {	
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("OK"))
}
