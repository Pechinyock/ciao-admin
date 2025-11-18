package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const inMemoryKey = "me-secret-key-ha-ha"

type JwtProvider struct{}

type TokenProvider interface {
	GenerateToken(string) (string, error)
	ValidateToken(string) error
}

func NewTokenProvider() (TokenProvider, error) {
	return &JwtProvider{}, nil
}

func (p *JwtProvider) ValidateToken(jwtStr string) error {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(inMemoryKey), nil
	})
	if err != nil || token.Valid {
		return err
	}
	return nil
}

func (p *JwtProvider) GenerateToken(userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userName,
		"exp":  time.Now().Add(25 * time.Hour).Unix(),
	})
	tokeStr, err := token.SignedString([]byte(inMemoryKey))
	if err != nil {
		return "", err
	}
	return tokeStr, nil
}
