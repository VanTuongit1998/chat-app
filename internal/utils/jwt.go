package utils

import (
	"errors"
	"time"

	"chat-app/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret string
	ttl    time.Duration
}

func NewJwtService(secret string, ttl time.Duration) *JwtService {
	return &JwtService{secret: secret, ttl: ttl}
}

func (s *JwtService) GenerateToken(user *model.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JwtService) ValidateToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return &model.User{Username: claims.Subject}, nil
}
