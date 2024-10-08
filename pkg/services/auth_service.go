package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	cfg AuthConfig
}

type AuthConfig struct {
	Salt      string
	Signature string
	TokenTTL  time.Duration
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{}
}
func (service *AuthService) hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum([]byte(service.cfg.Salt)))
}

func (service *AuthService) generateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.cfg.TokenTTL)),
	})

	return token.SignedString([]byte(service.cfg.Signature))
}

func (service *AuthService) ParseToken(header []string) (int64, error) {
	token, err := jwt.Parse(header[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(service.cfg.Signature), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, errors.New("invalid token")
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return int64(id), nil
}
