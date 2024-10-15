package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/implserv"
)

//go:generate mockgen -source=abstract.go -destination=mocks/mock.go

type ProjectService interface {
	Create(p dto.ProjectDTO, userId int64) (int64, error)
	GetById(id int64, userId int64) (dto.ProjectDTO, error)
	GetAll(userId int64) ([]dto.ProjectDTO, error)
	UpdateById(id int64, p dto.UpdateProjectDTO, userId int64) error
	DeleteById(id int64, userId int64) error
}

type AuthService interface {
	SignUp(su dto.SignUpDTO) (int64, error)
	SignIn(si dto.SignInDTO) (int64, error)
	HashPassword(password string) string
	GenerateTokens(id int64) (string, string, error)
	UpdateTokens(rt string) (string, string, error)
	ParseToken(input string) (int64, error)
}

type AbstractService struct {
	ProjectService
	AuthService
}

func NewService(repo *repositories.AbstractRepository, cfg *config.Config) *AbstractService {
	return &AbstractService{
		ProjectService: implserv.NewProjectService(repo.ProjectRepository),
		AuthService:    implserv.NewAuthService(repo.AuthRepository, cfg),
	}
}
