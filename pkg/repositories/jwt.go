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

func (repo *RefreshTokensRepository) Find(token string) error {
	return nil
}

func (repo *RefreshTokensRepository) Update(token string) error {
	return nil
}

func (repo *RefreshTokensRepository) Delete(token string) {

}
