package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
)

type ProjectService interface {
	Create(p dto.ProjectDTO) (int64, error)
	GetById(id int64) (dto.ProjectDTO, error)
	GetAll() ([]dto.ProjectDTO, error)
	UpdateById(id int64, p dto.UpdateProjectDTO) error
	DeleteById(id int64) error
}

type AbstractService struct {
	ProjectService
}

func NewService(repo *repositories.AbstractRepository) *AbstractService {
	return &AbstractService{
		ProjectService: NewProjectService(repo.ProjectRepository),
	}
}
