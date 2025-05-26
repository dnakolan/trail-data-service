package services

import (
	"context"
	"errors"
	"time"

	"github.com/dnakolan/trail-data-service/internal/config"
	"github.com/dnakolan/trail-data-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginService interface {
	GetSecretKey() []byte
	generateToken(username string) (string, error)
	Login(ctx context.Context, username string, password string) (string, error)
}

type loginService struct{}

func NewLoginService() *loginService {
	return &loginService{}
}

func (s *loginService) GetSecretKey() []byte {
	return []byte(config.SECRET_KEY)
}

func (s *loginService) generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(config.TOKEN_EXPIRATION_TIME)
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.TOKEN_ISSUER,
			Subject:   config.TOKEN_SUBJECT,
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.GetSecretKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *loginService) Login(ctx context.Context, username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}
	token, err := s.generateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
