package repositories

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type ProjectRepositoryImpl struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{db}
}

func (repo *ProjectRepositoryImpl) Add(p entity.Project) error {
	return nil
}

func (repo *ProjectRepositoryImpl) GetById(id int64) (entity.Project, error) {
	return entity.Project{}, nil
}

func (repo *ProjectRepositoryImpl) GetAll() ([]entity.Project, error) {
	return nil, nil
}

func (repo *ProjectRepositoryImpl) UpdateById(p entity.Project) error {
	return nil
}

func (repo *ProjectRepositoryImpl) DeleteById(id int64) error {
	return nil
}
