package services

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	repo repositories.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repositories.AuthRepository, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{
		repo: repo,
		cfg:  cfg,
	}
}

func (service *AuthServiceImpl) SignUp(su dto.SignUpDTO) (int64, error) {
	return service.repo.SignUp(entity.FromSignUpDTO(su, service.hashPassword(su.Password)))
}

func (service *AuthServiceImpl) SignIn(si dto.SignInDTO) (string, error) {
	id, err := service.repo.SignIn(si.Username, service.hashPassword(si.Password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	})

	return token.SignedString([]byte(service.cfg.Auth.Signature))
}

func (service *AuthServiceImpl) hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum([]byte(service.cfg.Auth.Salt)))
}
