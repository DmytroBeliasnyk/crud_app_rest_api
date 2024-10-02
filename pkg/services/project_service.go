package service

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories"
)

type ProjectServiceImpl struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{repo}
}

func (service *ProjectServiceImpl) Add(p dto.ProjectDTO) error {
	return nil
}

func (service *ProjectServiceImpl) GetById(id int64) (dto.ProjectDTO, error) {
	return dto.ProjectDTO{}, nil
}

func (service *ProjectServiceImpl) GetAll() ([]dto.ProjectDTO, error) {
	return nil, nil
}

func (service *ProjectServiceImpl) UpdateById(p dto.ProjectDTO) error {
	return nil
}

func (service *ProjectServiceImpl) DeleteById(id int64) error {
	return nil
}
