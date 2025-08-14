package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil }

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing error: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalide token")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != "chirpy" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, nil
}
