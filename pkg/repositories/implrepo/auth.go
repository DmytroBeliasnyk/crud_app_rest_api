package implrepo

import (
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (repo *UserRepositoryImpl) SignUp(u *entity.User) (int64, error) {
	var id int64
	if err := repo.db.QueryRow(`INSERT INTO users (name, email, username, password_hash)
	 							VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Name, u.Email, u.Username, u.PasswordHash).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *UserRepositoryImpl) SignIn(username, passwordHash string) (int64, error) {
	var id int64
	if err := repo.db.QueryRow("SELECT id FROM users WHERE username=$1 AND password_hash=$2",
		username, passwordHash).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *UserRepositoryImpl) CreateRefreshToken(userId int64, token string, expiresAt time.Time) error {
	if _, err := repo.db.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userId, token, expiresAt); err != nil {
		return err
	}

	return nil
}

func (repo *UserRepositoryImpl) FindRefreshToken(token string) (int64, time.Time, error) {
	var res struct {
		Id        int64     `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}

	if err := repo.db.Get(&res, "SELECT user_id, expires_at FROM tokens WHERE token=$1", token); err != nil {
		return 0, time.Time{}, err
	}

	return res.Id, res.ExpiresAt, nil
}

func (repo *UserRepositoryImpl) DeleteRefreshToken(token string) error {
	if _, err := repo.db.Exec("DELETE FROM tokens WHERE token=$1", token); err != nil {
		return err
	}

	return nil
}
