package implrepo

import (
	"fmt"
	"strings"

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

func (repo *ProjectRepositoryImpl) Create(p *entity.Project) (int64, error) {
	var id int64
	if err := repo.db.QueryRow("INSERT INTO projects (title, description, done) VALUES ($1, $2, $3) RETURNING id",
		p.Title, p.Description, p.Done).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ProjectRepositoryImpl) GetById(id int64) (entity.Project, error) {
	var project entity.Project
	if err := repo.db.Get(&project, "SELECT * FROM projects WHERE id=$1", id); err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (repo *ProjectRepositoryImpl) GetAll() (projects []entity.Project, err error) {
	if err = repo.db.Select(&projects, "SELECT * FROM projects"); err != nil {
		return nil, err
	}

	return projects, nil
}

func (repo *ProjectRepositoryImpl) UpdateById(id int64, input dto.UpdateProjectDTO) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	values := strings.Join(setValues, ", ")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE projects SET %s WHERE id=$%d", values, argId)
	if _, err := repo.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (repo *ProjectRepositoryImpl) DeleteById(id int64) error {
	if _, err := repo.db.Exec("DELETE FROM projects WHERE id=$1", id); err != nil {
		return err
	}

	return nil
}
