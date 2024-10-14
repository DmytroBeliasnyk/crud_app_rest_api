package implserv

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	mock_repositories "github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_Create(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockProjectRepository, p *entity.Project)

	cases := []struct {
		name         string
		input        dto.ProjectDTO
		inputUserId  int64
		mockBehavior mockBehavior
		expected     int64
		expectedErr  bool
	}{
		{
			name: "OK",
			input: dto.ProjectDTO{
				Title: "title",
			},
			inputUserId: 1,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, p *entity.Project) {
				s.EXPECT().Create(p).Return(int64(1), nil)
			},
			expected: 1,
		},
		{
			name:        "Emty field",
			inputUserId: 1,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, p *entity.Project) {
				s.EXPECT().Create(p).Return(int64(0), errors.New("some error"))
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockProjectRepository(ctrl)
			p := entity.FromDTO(c.input)
			p.UserId = c.inputUserId
			c.mockBehavior(repo, p)

			serv := NewProjectService(repo)
			userId, err := serv.Create(c.input, c.inputUserId)

			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userId, c.expected)
			}
		})
	}
}

func TestProjectService_GetById(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockProjectRepository, id int64, userId int64)

	cases := []struct {
		name         string
		inputId      int64
		inputUserId  int64
		mockBehavior mockBehavior
		expected     dto.ProjectDTO
		expectedErr  bool
	}{
		{
			name:        "OK",
			inputId:     1,
			inputUserId: 2,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, id, userId int64) {
				s.EXPECT().GetById(id, userId).Return(entity.Project{
					Id:     1,
					Title:  "title",
					Done:   false,
					UserId: 2,
				}, nil)
			},
			expected: dto.ProjectDTO{
				Title: "title",
			},
		},
		{
			name:        "Not found",
			inputId:     1,
			inputUserId: 2,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, id, userId int64) {
				s.EXPECT().GetById(id, userId).Return(entity.Project{}, sql.ErrNoRows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockProjectRepository(ctrl)
			c.mockBehavior(repo, c.inputId, c.inputUserId)

			serv := NewProjectService(repo)
			got, err := serv.GetById(c.inputId, c.inputUserId)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, c.expected)
			}
		})
	}
}

func TestProjectService_GetAll(t *testing.T) {
	expected := []dto.ProjectDTO{
		{Title: "title", Description: "description", Done: false},
		{Title: "title2", Description: "description2", Done: true},
	}

	mockBehavior := func(s *mock_repositories.MockProjectRepository, userId int64) {
		s.EXPECT().GetAll(userId).Return([]entity.Project{
			{Id: 1, Title: "title", Description: "description", Done: false, UserId: 1},
			{Id: 2, Title: "title2", Description: "description2", Done: true, UserId: 1},
		}, nil)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockProjectRepository(ctrl)
	mockBehavior(repo, 1)

	got, err := NewProjectService(repo).GetAll(1)

	assert.NoError(t, err)
	assert.Equal(t, got, expected)
}

func TestProjectService_UpdateById(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockProjectRepository,
		id int64, input dto.UpdateProjectDTO, userId int64)

	type args struct {
		input      dto.UpdateProjectDTO
		id, userId int64
	}

	cases := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expectedErr  bool
	}{
		{
			name: "OK",
			args: args{
				input: dto.UpdateProjectDTO{
					Title: stringPointer("New title"),
				},
				id:     1,
				userId: 2,
			},
			mockBehavior: func(s *mock_repositories.MockProjectRepository,
				id int64, input dto.UpdateProjectDTO, userId int64) {
				s.EXPECT().UpdateById(id, input, userId).Return(nil)
			},
		},
		{
			name: "Not found",
			args: args{
				input: dto.UpdateProjectDTO{
					Title: stringPointer("New title"),
				},
				id:     1,
				userId: 2,
			},
			mockBehavior: func(s *mock_repositories.MockProjectRepository,
				id int64, input dto.UpdateProjectDTO, userId int64) {
				s.EXPECT().UpdateById(id, input, userId).Return(errors.New("some error"))
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockProjectRepository(ctrl)
			c.mockBehavior(repo, c.args.id, c.args.input, c.args.userId)

			err := NewProjectService(repo).UpdateById(c.args.id, c.args.input, c.args.userId)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_DeleteById(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockProjectRepository, id int64, userId int64)

	cases := []struct {
		name         string
		inputId      int64
		inputUserId  int64
		mockBehavior mockBehavior
		expectedErr  bool
	}{
		{
			name:        "OK",
			inputId:     1,
			inputUserId: 2,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, id, userId int64) {
				s.EXPECT().DeleteById(id, userId).Return(nil)
			},
		},
		{
			name:        "Not found",
			inputId:     1,
			inputUserId: 2,
			mockBehavior: func(s *mock_repositories.MockProjectRepository, id, userId int64) {
				s.EXPECT().DeleteById(id, userId).Return(errors.New("some error"))
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockProjectRepository(ctrl)
			c.mockBehavior(repo, c.inputId, c.inputUserId)

			err := NewProjectService(repo).DeleteById(c.inputId, c.inputUserId)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPointer(str string) *string {
	return &str
}
