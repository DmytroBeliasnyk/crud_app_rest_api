package implrepo

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserRepository_SignUp(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewUserRepository(db)

	cases := []struct {
		name        string
		input       entity.User
		mock        func()
		expected    int64
		expectedErr bool
	}{
		{
			name: "OK",
			input: entity.User{
				Name:         "name",
				Email:        "aaa@bbb.ccc",
				Username:     "username",
				PasswordHash: "password",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("name", "aaa@bbb.ccc", "username", "password").
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name:  "emty fields",
			input: entity.User{},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("", "", "", "").
					WillReturnRows(rows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()

			got, err := repo.SignUp(&c.input)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, c.expected)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_SignIn(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewUserRepository(db)

	type args struct {
		username, passwordHash string
	}

	cases := []struct {
		name        string
		input       args
		mock        func()
		expected    int64
		expectedErr bool
	}{
		{
			name: "OK",
			input: args{
				username:     "username",
				passwordHash: "password",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs("username", "password").
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name:  "emty fields",
			input: args{},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs("", "").
					WillReturnRows(rows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()

			got, err := repo.SignIn(c.input.username, c.input.passwordHash)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, c.expected)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_CreateRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectExec("INSERT INTO tokens").
		WithArgs(1, "token", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)).
		WillReturnResult(driver.ResultNoRows)

	got := repo.CreateRefreshToken(1, "token", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local))

	assert.NoError(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewUserRepository(db)

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

			userId, expiresAt, err := repo.FindRefreshToken(c.arg)
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

func TestUserRepository_DeleteRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	assert.NoError(t, err, "an error '%s' was not expected when opening a stub database connection", err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectExec("DELETE FROM tokens").
		WithArgs("token").
		WillReturnResult(driver.ResultNoRows)

	err = repo.DeleteRefreshToken("token")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
