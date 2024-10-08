package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type RefreshTokensRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokensRepository {
	return &RefreshTokensRepository{db}
}

func (repo *RefreshTokensRepository) Create(userId int64, token string, expiresAt time.Time) error {
	if _, err := repo.db.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userId, token, expiresAt); err != nil {
		return err
	}

	return nil
}

func (repo *RefreshTokensRepository) Find(token string) (int64, time.Time, error) {
	var res struct {
		Id        int64     `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}

	if err := repo.db.Get(&res, "SELECT user_id, expires_at FROM tokens WHERE token=$1", token); err != nil {
		return 0, time.Time{}, err
	}

	return res.Id, res.ExpiresAt, nil
}

func (repo *RefreshTokensRepository) Delete(token string) error {
	if _, err := repo.db.Exec("DELETE FROM tokens WHERE token=$1", token); err != nil {
		return err
	}

	return nil
}
