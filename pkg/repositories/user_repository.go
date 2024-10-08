package repositories

import (
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (repo *UserRepositoryImpl) Create(u *entity.User) (int64, error) {
	var id int64
	if err := repo.db.QueryRow(`INSERT INTO users (name, email, username, password_hash)
	 							VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Name, u.Email, u.Username, u.PasswordHash).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *UserRepositoryImpl) Find(username, passwordHash string) (int64, error) {
	var id int64
	if err := repo.db.QueryRow("SELECT id FROM users WHERE username=$1 AND password_hash=$2",
		username, passwordHash).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
