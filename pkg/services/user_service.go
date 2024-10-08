package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
)

type UserServiceImpl struct {
	repo repositories.UserRepository
	auth *AuthService
}

func NewUserService(repo repositories.UserRepository, serv *AuthService) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
		auth: serv,
	}
}

func (service *UserServiceImpl) SignUp(su dto.SignUpDTO) (int64, error) {
	return service.repo.Create(entity.FromSignUpDTO(su, service.auth.hashPassword(su.Password)))
}

func (service *UserServiceImpl) SignIn(si dto.SignInDTO) (string, error) {
	id, err := service.repo.Find(si.Username, service.auth.hashPassword(si.Password))
	if err != nil {
		return "", err
	}

	return service.auth.generateToken(id)
}
