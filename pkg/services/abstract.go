package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/implserv"
)

type ProjectService interface {
	Create(p dto.ProjectDTO) (int64, error)
	GetById(id int64) (dto.ProjectDTO, error)
	GetAll() ([]dto.ProjectDTO, error)
	UpdateById(id int64, p dto.UpdateProjectDTO) error
	DeleteById(id int64) error
}

type UserService interface {
	SignUp(su dto.SignUpDTO) (int64, error)
	SignIn(si dto.SignInDTO) (string, string, error)
}

type AbstractService struct {
	ProjectService
	UserService
}

func NewService(repo *repositories.AbstractRepository, serv *implserv.AuthService) *AbstractService {
	return &AbstractService{
		ProjectService: implserv.NewProjectService(repo.ProjectRepository),
		UserService:    implserv.NewUserService(repo.UserRepository, serv),
	}
}
