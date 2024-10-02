package datasource

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository interface {
	Add(p entity.Project) error
	GetById(id int64) (entity.Project, error)
	GetAll() ([]entity.Project, error)
	UpdateById(p entity.Project) error
	DeleteById(id int64) error
}

type Repository struct {
	ProjectRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		ProjectRepository: NewProjectRepository(db),
	}
}
