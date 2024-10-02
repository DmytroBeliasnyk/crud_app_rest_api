package entity

import "github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"

type Project struct {
	Id          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Done        bool   `db:"done"`
}

func FromDTO(dto dto.ProjectDTO) *Project {
	return &Project{
		Title:       dto.Title,
		Description: dto.Description,
		Done:        dto.Done,
	}
}

func (p *Project) ToDTO() *dto.ProjectDTO {
	return &dto.ProjectDTO{
		Title:       p.Title,
		Description: p.Description,
		Done:        p.Done,
	}
}
