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
	if err := repo.db.QueryRow(`INSERT INTO projects (title, description, done, user_id)
								 VALUES ($1, $2, $3, $4) RETURNING id`,
		p.Title, p.Description, p.Done, p.UserId).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ProjectRepositoryImpl) GetById(id int64, userId int64) (entity.Project, error) {
	var project entity.Project
	if err := repo.db.Get(&project, "SELECT * FROM projects WHERE id=$1 AND user_id=$2", id, userId); err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (repo *ProjectRepositoryImpl) GetAll(userId int64) (projects []entity.Project, err error) {
	if err = repo.db.Select(&projects, "SELECT * FROM projects WHERE user_id=$1", userId); err != nil {
		return nil, err
	}

	return projects, nil
}

func (repo *ProjectRepositoryImpl) UpdateById(id int64, input dto.UpdateProjectDTO, userId int64) error {
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
	args = append(args, id, userId)

	query := fmt.Sprintf("UPDATE projects SET %s WHERE id=$%d AND user_id=$%d", values, argId, argId+1)
	if _, err := repo.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (repo *ProjectRepositoryImpl) DeleteById(id int64, userId int64) error {
	if _, err := repo.db.Exec("DELETE FROM projects WHERE id=$1 AND user_id=$2", id, userId); err != nil {
		return err
	}

	return nil
}
