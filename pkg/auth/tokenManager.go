package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

//go:generate mockgen -source=tokenManager.go -destination=../../internal/tests/unitTests/serviceTests/mocks/mockTokenManager.go --package=mocks

// ITokenManager provides logic for JWT & Refresh tokens generation and parsing.
type ITokenManager interface {
	NewJWT(userID uuid.UUID, role string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, string, error)
	NewRefreshToken() (string, error)
}

type tokenClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type TokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) (ITokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &TokenManager{signingKey: signingKey}, nil
}

func (m *TokenManager) NewJWT(readerID uuid.UUID, role string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   readerID.String(),
		},
		Role: role,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) Parse(accessToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", "", fmt.Errorf("error get user claims from token")
	}

	return claims.RegisteredClaims.Subject, claims.Role, nil
}

func (m *TokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
