package jwt

import (
	"fmt"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenService struct{ secret []byte }

func NewTokenService(secret string) *TokenService { return &TokenService{secret: []byte(secret)} }

type CustomClaims struct {
	Username string `json:"username"`
	jwtlib.RegisteredClaims
}

func (s *TokenService) GenerateToken(username, sessionID string) (string, error) {
	now := time.Now()
	claims := CustomClaims{Username: username, RegisteredClaims: jwtlib.RegisteredClaims{ID: sessionID, Subject: username, IssuedAt: jwtlib.NewNumericDate(now), ExpiresAt: jwtlib.NewNumericDate(now.Add(365 * 24 * time.Hour)), NotBefore: jwtlib.NewNumericDate(now)}}
	value, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return value, nil
}
