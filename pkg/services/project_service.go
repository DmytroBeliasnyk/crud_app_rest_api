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
	project, err := service.repo.GetById(id)

	return *project.ToDTO(), err
}

func (service *ProjectServiceImpl) GetAll() ([]dto.ProjectDTO, error) {
	projects, err := service.repo.GetAll()

	dtos := make([]dto.ProjectDTO, len(projects))
	for i, p := range projects {
		dtos[i] = *p.ToDTO()
	}

	return dtos, err
}

func (service *ProjectServiceImpl) UpdateById(id int64, input dto.UpdateProjectDTO) error {
	return service.repo.UpdateById(id, input)
}

func (service *ProjectServiceImpl) DeleteById(id int64) error {
	return service.repo.DeleteById(id)
}
