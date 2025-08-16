package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/simoncdn/http-server/internal/auth"
	"github.com/simoncdn/http-server/internal/config"
	"github.com/simoncdn/http-server/internal/database"
)

type UserHandler struct {
	cfg *config.Config
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func NewUserHanlder(cfg *config.Config) *UserHandler {
	userHandler := &UserHandler{
		cfg: cfg,
	}
	return userHandler
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type UserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var userRequest UserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Fatal("couldn't decode user request: %w", err)
	}

	email := userRequest.Email
	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		log.Fatal("hash password error: %w", err)
	}

	newUser := database.CreateUserParams {
		Email: email,
		HashedPassword: hashedPassword,
	}

	user, err := u.cfg.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		log.Fatal("couldn't create a new user: %w", err)
	}

	userFormatted := MapUserToResponse(user)

	data, err := json.Marshal(userFormatted)
	if err != nil {
		log.Fatal("error on marshalling user: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(data))
}

func MapUserToResponse(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
