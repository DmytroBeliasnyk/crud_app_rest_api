package implrepo

import (
	"slices"
	"testing"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	cases := []struct {
		name        string
		project     entity.Project
		mock        func()
		expected    int64
		expectedErr bool
	}{
		{
			name: "OK",
			project: entity.Project{
				Title:  "title",
				UserId: 1,
			},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO projects").
					WithArgs("title", "", false, 1).
					WillReturnRows(rows)
			},
			expected: 1,
		},
		{
			name:    "emty field",
			project: entity.Project{},
			mock: func() {
				mock.ExpectQuery("INSERT INTO projects").
					WithArgs("", "", false, 0)
			},
			expected:    1,
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()
			got, err := repo.Create(&c.project)
			if (err != nil) != c.expectedErr {
				t.Errorf("got error = %s, expected error %t", err, c.expectedErr)
				return
			}
			if err == nil && got != c.expected {
				t.Errorf("got = %d, expected %d", got, c.expected)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	type args struct {
		id     int64
		userId int64
	}

	cases := []struct {
		name        string
		args        args
		mock        func()
		expected    entity.Project
		expectedErr bool
	}{
		{
			name: "OK",
			args: args{1, 2},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"id", "title", "description", "done", "user_id"}).
					AddRow(1, "title", "description", false, 2)
				mock.ExpectQuery("SELECT (.+) FROM projects").
					WithArgs(1, 2).
					WillReturnRows(rows)
			},
			expected: entity.Project{
				Id:          1,
				Title:       "title",
				Description: "description",
				Done:        false,
				UserId:      2,
			},
		},
		{
			name: "Not found",
			args: args{1, 2},
			mock: func() {
				rows := sqlxmock.NewRows([]string{"id", "title", "description", "done", "user_id"})
				mock.ExpectQuery("SELECT (.+) FROM projects").
					WithArgs(1, 2).
					WillReturnRows(rows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mock()

			got, err := repo.GetById(c.args.id, c.args.userId)
			if (err != nil) != c.expectedErr {
				t.Fatalf("got err = %s, expected error = %t", err, c.expectedErr)
				return
			}
			if err == nil && got != c.expected {
				t.Fatalf("got err = %v, expected error = %v", got, c.expected)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	arg := int64(1)
	expected := []entity.Project{
		{
			Id:          1,
			Title:       "title",
			Description: "description",
			UserId:      arg,
		},
	}
	rows := sqlxmock.NewRows([]string{"id", "title", "description", "done", "user_id"}).
		AddRow(1, "title", "description", false, arg)

	mock.ExpectQuery("SELECT (.+) FROM projects").
		WithArgs(arg).
		WillReturnRows(rows)

	got, err := repo.GetAll(arg)
	if err != nil {
		t.Fatalf("got err = %s", err)
	}
	if !slices.Equal(got, expected) {
		t.Fatalf("got err = %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateById(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	updateDone := true
	input := dto.UpdateProjectDTO{
		Done: &updateDone,
	}
	mock.ExpectExec("UPDATE projects SET").
		WithArgs(updateDone, 1, 2).
		WillReturnResult(sqlxmock.NewResult(0, 1))

	if got := repo.UpdateById(1, input, 2); got != nil {
		t.Errorf("got error: %s", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteById(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewProjectRepository(db)

	mock.ExpectExec("DELETE FROM projects").
		WithArgs(1, 2).
		WillReturnResult(sqlxmock.NewResult(0, 1))

	if got := repo.DeleteById(1, 2); got != nil {
		t.Errorf("got error: %s", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
