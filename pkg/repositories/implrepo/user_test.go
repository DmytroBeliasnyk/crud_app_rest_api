package implrepo

import (
	"testing"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserRepository_Create(t *testing.T) {
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

			got, err := repo.Create(&c.input)
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

func TestUserRepository_Find(t *testing.T) {
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

			got, err := repo.Find(c.input.username, c.input.passwordHash)
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
