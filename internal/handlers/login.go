package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/simoncdn/http-server/internal/auth"
	"github.com/simoncdn/http-server/internal/config"
)

type LoginHander struct {
	cfg *config.Config
}

func NewLoginHandler(config *config.Config) *LoginHander {
	newLoginHandler := &LoginHander{
		cfg: config,
	}

	return newLoginHandler
}

func (h *LoginHander) Login(w http.ResponseWriter, r *http.Request) {
	const default_expiration = 3600
	type Parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)

	parameters := Parameters{
		ExpiresInSeconds: default_expiration,
	}
	err := decoder.Decode(&parameters)
	if err != nil {
		fmt.Println("decode parameters error: ", err)
		return
	}

	dbUser, err := h.cfg.DB.GetUserByEmail(r.Context(), parameters.Email)
	if err != nil {
		fmt.Println("couldn't get the current user: ", err)
		return
	}

	err = auth.CheckPasswordHash(parameters.Password, dbUser.HashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Hour
	if parameters.ExpiresInSeconds > 0 && parameters.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(parameters.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, h.cfg.JWTSecret, expirationTime)
	if err != nil {
		log.Fatal("couldn't generate a JWT", err)
	}

	user := response{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: accessToken,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("couldn't marshal user: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
