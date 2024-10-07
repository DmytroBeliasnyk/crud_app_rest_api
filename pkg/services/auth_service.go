package services

import (
	"crypto/sha256"
	"fmt"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
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

func (service *AuthServiceImpl) hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum([]byte(service.cfg.Auth.Salt)))
}
