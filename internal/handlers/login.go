package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	type Parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	parameters := Parameters{}
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

	user := MapUserToResponse(dbUser)

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("couldn't marshal user: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
