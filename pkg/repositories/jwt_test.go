package repositories

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRefreshTokenRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewRefreshTokenRepository(db)

	mock.ExpectExec("INSERT INTO tokens").
		WithArgs(1, "token", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)).
		WillReturnResult(driver.ResultNoRows)

	got := repo.Create(1, "token", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local))

	assert.NoError(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRefreshTokenRepository_Find(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewRefreshTokenRepository(db)

	type result struct {
		userId    int64
		expiresAt time.Time
	}
	cases := []struct {
		name        string
		arg         string
		mock        func()
		expected    result
		expectedErr bool
	}{
		{
			name: "OK",
			arg:  "token",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "expires_at"}).
					AddRow(1, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local))

				mock.ExpectQuery("SELECT (.+) FROM tokens").
					WithArgs("token").
					WillReturnRows(rows)
			},
			expected: result{
				userId:    1,
				expiresAt: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "Not found",
			arg:  "not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "expires_at"})

				mock.ExpectQuery("SELECT (.+) FROM tokens").
					WithArgs("not found").
					WillReturnRows(rows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()

			userId, expiresAt, err := repo.Find(c.arg)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userId, c.expected.userId)
				assert.Equal(t, expiresAt, c.expected.expiresAt)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshTokenRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewRefreshTokenRepository(db)

	mock.ExpectExec("DELETE FROM tokens").
		WithArgs("token").
		WillReturnResult(driver.ResultNoRows)

	err = repo.Delete("token")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
