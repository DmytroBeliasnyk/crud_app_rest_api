package repositories

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type ProjectRepositoryImpl struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{db}
}

func (repo *ProjectRepositoryImpl) Create(p entity.Project) (int64, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, tx.Rollback()
	}

	var id int64
	if err = tx.QueryRow("INSERT INTO projects (title, description, done) VALUES ($1, $2, $3) RETURNING id",
		p.Title, p.Description, p.Done).
		Scan(&id); err != nil {
		return 0, tx.Rollback()
	}

	return id, tx.Commit()
}

func (repo *ProjectRepositoryImpl) GetById(id int64) (dto.ProjectDTO, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return dto.ProjectDTO{}, tx.Rollback()
	}

	var project entity.Project
	if err = tx.Get(&project, "SELECT * FROM projects WHERE id=$1", id); err != nil {
		return dto.ProjectDTO{}, tx.Rollback()
	}

	return *project.ToDTO(), tx.Commit()
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
