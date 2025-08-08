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
	mux.Handle("/", fileServer)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("serving on port: %s", port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error on serving: %w", err)
	}

	return nil
}
