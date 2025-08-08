package main

import (
	"fmt"
	"net/http"
)

const (
	port = "8080"
)

func server() error {
	mux := http.NewServeMux()

	server := http.Server {
		Addr: ":" + port,
		Handler: mux,
	}

	fmt.Printf("serving on port: %s", port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error on serving: %w", err)
	}

	return nil
}
