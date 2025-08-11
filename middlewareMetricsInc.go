package main

import (
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	handler := func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	}
	handlerFunc := http.HandlerFunc(handler)

	return handlerFunc
}
