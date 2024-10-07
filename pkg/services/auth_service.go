package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
)

type AuthServiceImpl struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{repo}
}

func (service *AuthServiceImpl) SignUp(su dto.SignUpDTO) (int64, error) {
	return service.repo.SignUp(entity.FromSignUpDTO(su, su.Password))
}
