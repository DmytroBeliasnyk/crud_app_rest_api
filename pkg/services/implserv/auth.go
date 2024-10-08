package implserv

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo *repositories.RefreshTokensRepository
	cfg  authConfig
}

type authConfig struct {
	salt      string
	signature string
	jwt       time.Duration
	refresh   time.Duration
}

func NewAuthService(repo *repositories.RefreshTokensRepository, config *config.Config) *AuthService {
	auth := config.Auth
	return &AuthService{
		repo: repo,
		cfg: authConfig{
			salt:      auth.Salt,
			signature: auth.Signature,
			jwt:       auth.JWT,
			refresh:   auth.JWT,
		},
	}
}

func (service *AuthService) HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum([]byte(service.cfg.salt)))
}

func (service *AuthService) GenerateTokens(id int64) (string, string, error) {
	jwtt := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.cfg.jwt)),
	})
	jt, err := jwtt.SignedString([]byte(service.cfg.signature))
	if err != nil {
		return "", "", err
	}

	token := make([]byte, 32)
	if _, err = rand.Read(token); err != nil {
		return "", "", err
	}

	rt := fmt.Sprintf("%x", token)
	if err = service.repo.Create(id, rt, time.Now().Add(service.cfg.refresh)); err != nil {
		return "", "", err
	}

	return jt, rt, nil
}

func (service *AuthService) UpdateTokens(rt string) (string, string, error) {
	id, expiresAt, err := service.repo.Find(rt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("invalid refresh token")
		}
		return "", "", err
	}

	if time.Now().After(expiresAt) {
		if err := service.repo.Delete(rt); err != nil {
			return "", "", err
		}
		return "", "", errors.New("token is expired")
	}

	if err := service.repo.Delete(rt); err != nil {
		return "", "", err
	}

	return service.GenerateTokens(id)
}

func (service *AuthService) ParseToken(input string) (int64, error) {
	token, err := jwt.Parse(input, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(service.cfg.signature), nil
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
