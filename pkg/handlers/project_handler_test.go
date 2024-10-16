package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	mock_handlers "github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/handlers/mocks"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	mock_services "github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_create(t *testing.T) {
	type serviceBehavior func(s *mock_services.MockProjectService, input dto.ProjectDTO, userId int64)
	type cacheBehavior func(s *mock_handlers.MockCache, userId int64)
	cases := []struct {
		name                string
		body                string
		input               dto.ProjectDTO
		userId              int64
		serviceBehavior     serviceBehavior
		cacheBehavior       cacheBehavior
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name: "OK",
			body: `{"title":"title"}`,
			input: dto.ProjectDTO{
				Title: "title",
			},
			userId: 1,
			serviceBehavior: func(s *mock_services.MockProjectService, input dto.ProjectDTO, userId int64) {
				s.EXPECT().Create(input, userId).Return(userId, nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, userId int64) {
				s.EXPECT().Delete(fmt.Sprintf("all%d", userId))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:                "Invalid input body",
			body:                `{"":""}`,
			input:               dto.ProjectDTO{},
			userId:              1,
			serviceBehavior:     func(s *mock_services.MockProjectService, input dto.ProjectDTO, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, userId int64) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "User unauthorized",
			body: `{"title":"title"}`,
			input: dto.ProjectDTO{
				Title: "title",
			},
			userId:              0,
			serviceBehavior:     func(s *mock_services.MockProjectService, input dto.ProjectDTO, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, userId int64) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name: "Service failed",
			body: `{"title":"title"}`,
			input: dto.ProjectDTO{
				Title: "title",
			},
			userId: 1,
			serviceBehavior: func(s *mock_services.MockProjectService, input dto.ProjectDTO, userId int64) {
				s.EXPECT().Create(input, userId).Return(int64(0), errors.New("some error"))
			},
			cacheBehavior:       func(s *mock_handlers.MockCache, userId int64) {},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServ := mock_services.NewMockProjectService(ctrl)
			c.serviceBehavior(mockServ, c.input, c.userId)

			mockCache := mock_handlers.NewMockCache(ctrl)
			c.cacheBehavior(mockCache, c.userId)

			serv := services.AbstractService{ProjectService: mockServ}
			h := Handler{
				service: &serv,
				cache:   mockCache,
			}

			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("user_id", c.userId)
			})
			r.POST("/create", h.create)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(c.body))

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				msg, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, msg)
			} else {
				assert.Equal(t, rec.Body.String(), `{"id":1}`)
			}
		})
	}
}

func TestHandler_getById(t *testing.T) {
	type mockService func(s *mock_services.MockProjectService, projectId, userId int64)
	type mockCache func(s *mock_handlers.MockCache, projectId, userId int64)

	cases := []struct {
		name                string
		projectId           int64
		userId              int64
		serviceBehavior     mockService
		cacheBehavior       mockCache
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name:      "OK",
			projectId: 1,
			userId:    2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {
				s.EXPECT().GetById(projectId, userId).Return(dto.ProjectDTO{Title: "title"}, nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				cache := fmt.Sprintf("%d%d", projectId, userId)

				s.EXPECT().Get(cache).Return(gomock.Any(), errors.New("some error"))
				s.EXPECT().Set(cache, dto.ProjectDTO{Title: "title"}, gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "OK from cache",
			projectId:       1,
			userId:          2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				s.EXPECT().Get(fmt.Sprintf("%d%d", projectId, userId)).
					Return(dto.ProjectDTO{Title: "title"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "Invalid query",
			userId:              2,
			serviceBehavior:     func(s *mock_services.MockProjectService, projectId, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:                "User unauthorized",
			projectId:           1,
			serviceBehavior:     func(s *mock_services.MockProjectService, projectId, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name:      "Service failed",
			projectId: 1,
			userId:    2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {
				s.EXPECT().GetById(projectId, userId).Return(dto.ProjectDTO{}, errors.New("some error"))
			},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				s.EXPECT().Get(fmt.Sprintf("%d%d", projectId, userId)).
					Return(gomock.Any, errors.New("some error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
		{
			name:      "Set cache failed",
			projectId: 1,
			userId:    2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {
				s.EXPECT().GetById(projectId, userId).
					Return(dto.ProjectDTO{Title: "title"}, nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				cache := fmt.Sprintf("%d%d", projectId, userId)

				s.EXPECT().Get(cache).Return(gomock.Any, errors.New("some error"))
				s.EXPECT().Set(cache, dto.ProjectDTO{Title: "title"}, gomock.Any()).Return(errors.New("some error"))
				s.EXPECT().Delete(cache)
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_services.NewMockProjectService(ctrl)
			c.serviceBehavior(serviceMock, c.projectId, c.userId)

			cacheMock := mock_handlers.NewMockCache(ctrl)
			c.cacheBehavior(cacheMock, c.projectId, c.userId)

			serv := services.AbstractService{ProjectService: serviceMock}
			h := Handler{
				service: &serv,
				cache:   cacheMock,
			}

			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("user_id", c.userId)
			})
			target := "/get-by-id"
			r.GET(target, h.getById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("%s?id=%d", target, c.projectId), nil)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				msg, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, msg)
			} else {
				assert.Equal(t, rec.Body.String(), `{"title":"title","description":"","done":false}`)
			}
		})
	}
}

func TestHandler_getAll(t *testing.T) {
	type mockService func(s *mock_services.MockProjectService, userId int64)
	type mockCache func(s *mock_handlers.MockCache, userId int64)

	cases := []struct {
		name                string
		userId              int64
		serviceBehavior     mockService
		cacheBehavior       mockCache
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name:   "OK",
			userId: 2,
			serviceBehavior: func(s *mock_services.MockProjectService, userId int64) {
				s.EXPECT().GetAll(userId).Return([]dto.ProjectDTO{{Title: "title"}}, nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, userId int64) {
				cache := fmt.Sprintf("all%d", userId)

				s.EXPECT().Get(cache).Return(gomock.Any(), errors.New("some error"))
				s.EXPECT().Set(cache, []dto.ProjectDTO{{Title: "title"}}, gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "OK from cache",
			userId:          2,
			serviceBehavior: func(s *mock_services.MockProjectService, userId int64) {},
			cacheBehavior: func(s *mock_handlers.MockCache, userId int64) {
				s.EXPECT().Get(fmt.Sprintf("all%d", userId)).
					Return([]dto.ProjectDTO{{Title: "title"}}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "User unauthorized",
			serviceBehavior:     func(s *mock_services.MockProjectService, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, userId int64) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name:   "Service failed",
			userId: 2,
			serviceBehavior: func(s *mock_services.MockProjectService, userId int64) {
				s.EXPECT().GetAll(userId).Return([]dto.ProjectDTO{}, errors.New("some error"))
			},
			cacheBehavior: func(s *mock_handlers.MockCache, userId int64) {
				s.EXPECT().Get(fmt.Sprintf("all%d", userId)).
					Return(gomock.Any, errors.New("some error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
		{
			name:   "Set cache failed",
			userId: 2,
			serviceBehavior: func(s *mock_services.MockProjectService, userId int64) {
				s.EXPECT().GetAll(userId).Return([]dto.ProjectDTO{{Title: "title"}}, nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, userId int64) {
				cache := fmt.Sprintf("all%d", userId)

				s.EXPECT().Get(cache).Return(gomock.Any, errors.New("some error"))
				s.EXPECT().Set(cache, []dto.ProjectDTO{{Title: "title"}}, gomock.Any()).Return(errors.New("some error"))
				s.EXPECT().Delete(cache)
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_services.NewMockProjectService(ctrl)
			c.serviceBehavior(serviceMock, c.userId)

			cacheMock := mock_handlers.NewMockCache(ctrl)
			c.cacheBehavior(cacheMock, c.userId)

			serv := services.AbstractService{ProjectService: serviceMock}
			h := Handler{
				service: &serv,
				cache:   cacheMock,
			}

			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("user_id", c.userId)
			})
			target := "/get-all"
			r.GET(target, h.getAll)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", target, nil)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				msg, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, msg)
			} else {
				assert.Equal(t, rec.Body.String(), `[{"title":"title","description":"","done":false}]`)
			}
		})
	}
}

func TestHandler_updateById(t *testing.T) {
	type mockService func(s *mock_services.MockProjectService, projectId int64,
		input dto.UpdateProjectDTO, userId int64)
	type mockCache func(s *mock_handlers.MockCache, projectId, userId int64)

	cases := []struct {
		name                string
		projectId           int64
		userId              int64
		body                string
		input               dto.UpdateProjectDTO
		serviceBehavior     mockService
		cacheBehavior       mockCache
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name:      "OK",
			projectId: 1,
			userId:    2,
			body:      `{"done":true}`,
			input:     dto.UpdateProjectDTO{Done: boolPointer(true)},
			serviceBehavior: func(s *mock_services.MockProjectService, projectId int64,
				input dto.UpdateProjectDTO, userId int64) {
				s.EXPECT().UpdateById(projectId, input, userId).Return(nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				s.EXPECT().Delete(fmt.Sprintf("%d%d", projectId, userId))
				s.EXPECT().Delete(fmt.Sprintf("all%d", userId))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid query",
			userId: 2,
			body:   `{"done":true}`,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId int64,
				input dto.UpdateProjectDTO, userId int64) {
			},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:      "User unauthorized",
			projectId: 1,
			body:      `{"done":true}`,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId int64,
				input dto.UpdateProjectDTO, userId int64) {
			},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name:      "Service failed",
			projectId: 1,
			userId:    2,
			input:     dto.UpdateProjectDTO{Done: boolPointer(true)},
			body:      `{"done":true}`,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId int64,
				input dto.UpdateProjectDTO, userId int64) {
				s.EXPECT().UpdateById(projectId, input, userId).Return(errors.New("some error"))
			},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_services.NewMockProjectService(ctrl)
			c.serviceBehavior(serviceMock, c.projectId, c.input, c.userId)

			cacheMock := mock_handlers.NewMockCache(ctrl)
			c.cacheBehavior(cacheMock, c.projectId, c.userId)

			serv := services.AbstractService{ProjectService: serviceMock}
			h := Handler{
				service: &serv,
				cache:   cacheMock,
			}

			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("user_id", c.userId)
			})
			target := "/update"
			r.POST(target, h.updateById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("%s?id=%d", target, c.projectId),
				bytes.NewBufferString(c.body))

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				msg, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, msg)
			} else {
				assert.Equal(t, rec.Body.String(), `{"message":"ok"}`)
			}
		})
	}
}

func TestHandler_deleteById(t *testing.T) {
	type mockService func(s *mock_services.MockProjectService, projectId, userId int64)
	type mockCache func(s *mock_handlers.MockCache, projectId, userId int64)

	cases := []struct {
		name                string
		projectId           int64
		userId              int64
		serviceBehavior     mockService
		cacheBehavior       mockCache
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name:      "OK",
			projectId: 1,
			userId:    2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {
				s.EXPECT().DeleteById(projectId, userId).Return(nil)
			},
			cacheBehavior: func(s *mock_handlers.MockCache, projectId, userId int64) {
				s.EXPECT().Delete(fmt.Sprintf("%d%d", projectId, userId))
				s.EXPECT().Delete(fmt.Sprintf("all%d", userId))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "Invalid query",
			userId:              2,
			serviceBehavior:     func(s *mock_services.MockProjectService, projectId, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:                "User unauthorized",
			projectId:           1,
			serviceBehavior:     func(s *mock_services.MockProjectService, projectId, userId int64) {},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name:      "Service failed",
			projectId: 1,
			userId:    2,
			serviceBehavior: func(s *mock_services.MockProjectService, projectId, userId int64) {
				s.EXPECT().DeleteById(projectId, userId).Return(errors.New("some error"))
			},
			cacheBehavior:       func(s *mock_handlers.MockCache, projectId, userId int64) {},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_services.NewMockProjectService(ctrl)
			c.serviceBehavior(serviceMock, c.projectId, c.userId)

			cacheMock := mock_handlers.NewMockCache(ctrl)
			c.cacheBehavior(cacheMock, c.projectId, c.userId)

			serv := services.AbstractService{ProjectService: serviceMock}
			h := Handler{
				service: &serv,
				cache:   cacheMock,
			}

			r := gin.New()
			r.Use(func(ctx *gin.Context) {
				ctx.Set("user_id", c.userId)
			})
			target := "/delete"
			r.DELETE(target, h.deleteById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("%s?id=%d", target, c.projectId), nil)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				msg, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, msg)
			} else {
				assert.Equal(t, rec.Body.String(), `{"message":"ok"}`)
			}
		})
	}
}

func boolPointer(b bool) *bool {
	return &b
}
