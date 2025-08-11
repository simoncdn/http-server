package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	cleanedBody := getCleanedBody(parameters.Body)
	respondWithJson(w, http.StatusOK, cleanedBody)
}

func getCleanedBody(message string) map[string]string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(message)
	for i, word := range words {
		loweredWord := strings.ToLower(word)

		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	cleanedMessage := strings.Join(words, " ")
	return map[string]string{"cleaned_body": cleanedMessage}
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
