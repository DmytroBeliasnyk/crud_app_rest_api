package services

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
)

type ProjectServiceImpl struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{repo}
}

func (service *ProjectServiceImpl) Create(p dto.ProjectDTO) (int64, error) {
	return service.repo.Create(*entity.FromDTO(p))
}

func (service *ProjectServiceImpl) GetById(id int64) (dto.ProjectDTO, error) {
	return service.repo.GetById(id)
}

func (service *ProjectServiceImpl) GetAll() ([]dto.ProjectDTO, error) {
	return nil, nil
}

func (service *ProjectServiceImpl) UpdateById(id int64, p dto.UpdateProjectDTO) error {
	return nil
}

func (service *ProjectServiceImpl) DeleteById(id int64) error {
	return nil
}
