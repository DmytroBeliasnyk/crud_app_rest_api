package repositories

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository interface {
	Create(p entity.Project) (int64, error)
	GetById(id int64) (entity.Project, error)
	GetAll() ([]entity.Project, error)
	UpdateById(p entity.Project) error
	DeleteById(id int64) error
}

type AbstractRepository struct {
	ProjectRepository
}

func NewRepository(db *sqlx.DB) *AbstractRepository {
	return &AbstractRepository{
		ProjectRepository: NewProjectRepository(db),
	}
}
