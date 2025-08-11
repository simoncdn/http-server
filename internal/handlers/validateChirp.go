package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Parameters struct {
	Body string `json:"body"`
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	parameters := Parameters{}

	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, map[string]string{"error": "Something went wrong"})
		return
	}

	if len(parameters.Body) > 140 {
		respondWithJson(w, http.StatusBadRequest, map[string]string{"error": "Chirp is too long"})
		return
	}

	respondWithJson(w, http.StatusOK, map[string]bool{"isValid": true})
}

func respondWithJson(w http.ResponseWriter, statusCode int, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
