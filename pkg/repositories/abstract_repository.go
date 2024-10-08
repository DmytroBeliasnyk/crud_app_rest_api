package repositories

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository interface {
	Create(p *entity.Project) (int64, error)
	GetById(id int64) (entity.Project, error)
	GetAll() ([]entity.Project, error)
	UpdateById(id int64, input dto.UpdateProjectDTO) error
	DeleteById(id int64) error
}

type UserRepository interface {
	Create(u *entity.User) (int64, error)
	Find(username, passwordHash string) (int64, error)
}

type AbstractRepository struct {
	ProjectRepository
	UserRepository
}

func NewRepository(db *sqlx.DB) *AbstractRepository {
	return &AbstractRepository{
		ProjectRepository: NewProjectRepository(db),
		UserRepository:    NewUserRepository(db),
	}
}
