package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/simoncdn/http-server/internal/config"
	"github.com/simoncdn/http-server/internal/database"
)

type ChirpHandler struct {
	cfg *config.Config
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func NewChirpHanler(cfg *config.Config) *ChirpHandler {
	newChirpHanlder := &ChirpHandler{
		cfg: cfg,
	}

	return newChirpHanlder
}

func (h *ChirpHandler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

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
	newChirp := database.CreateChirpParams{
		UserID: parameters.UserId,
		Body:   cleanedBody,
	}

	chirp, err := h.cfg.DB.CreateChirp(r.Context(), newChirp)
	if err != nil {
		fmt.Println("error on creating chirp:", err)
		return
	}

	formattedChirp := mapChirpToResponse(chirp)
	marchalledChirp, err := json.Marshal(formattedChirp)
	if err != nil {
		fmt.Println("error on creating chirp:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(marchalledChirp)
}

func (h *ChirpHandler) GetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := h.cfg.DB.GetChirps(r.Context())
	if err != nil {
		fmt.Println("couldn't fetch chirps: ", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		fmt.Println("couldn't marshal chirps: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func (h *ChirpHandler) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		fmt.Println("couldn't parse chirpId into uuid", err)
		return
	}

	fmt.Printf("chirpId: %s | chirpUUID: %s", chirpId, chirpUUID)

	chirp, err := h.cfg.DB.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	formattedChirp := mapChirpToResponse(chirp)
	data, err := json.Marshal(formattedChirp)
	if err != nil {
		fmt.Println("couldn't marshal chirps: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func mapChirpToResponse(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func getCleanedBody(message string) string {
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
	return cleanedMessage
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
