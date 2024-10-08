package implserv

import (
	"crypto/sha256"
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
	return &AuthService{
		repo: repo,
		cfg: authConfig{
			salt:      config.Auth.Salt,
			signature: config.Auth.Signature,
			jwt:       config.Auth.JWT,
			refresh:   config.Auth.Refresh,
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

	rwtt := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.cfg.refresh)),
	})
	rt, err := rwtt.SignedString([]byte(service.cfg.signature))
	if err != nil {
		return "", "", err
	}

	return jt, rt, nil
}

func (service *AuthService) ParseToken(header []string) (int64, error) {
	token, err := jwt.Parse(header[1], func(t *jwt.Token) (interface{}, error) {
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
