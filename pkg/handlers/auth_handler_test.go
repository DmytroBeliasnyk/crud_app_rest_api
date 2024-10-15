package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services"
	mock_services "github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/services/mocks"
	"github.com/DmytroBeliasnyk/in_memory_cache/memory"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthService, input dto.SignUpDTO)

	cases := []struct {
		name                string
		body                string
		input               dto.SignUpDTO
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name: "OK",
			body: `{"name":"Name","email":"aaa@bbb.ccc",
				"username":"username","password":"password"}`,
			input: dto.SignUpDTO{
				Name:     "Name",
				Email:    "aaa@bbb.ccc",
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignUpDTO) {
				s.EXPECT().SignUp(input).Return(int64(1), nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:                "Invalid body",
			body:                "invalid_body",
			input:               dto.SignUpDTO{},
			mockBehavior:        func(s *mock_services.MockAuthService, input dto.SignUpDTO) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "Invalid args in body",
			body: `{"name":"Name","email":"invalid_email",
				"username":"username","password":"password"}`,
			input:               dto.SignUpDTO{},
			mockBehavior:        func(s *mock_services.MockAuthService, input dto.SignUpDTO) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "Service failed",
			body: `{"name":"Name","email":"aaa@bbb.ccc",
				"username":"username","password":"password"}`,
			input: dto.SignUpDTO{
				Name:     "Name",
				Email:    "aaa@bbb.ccc",
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignUpDTO) {
				s.EXPECT().SignUp(input).Return(int64(0), errors.New("some error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_services.NewMockAuthService(ctrl)
			c.mockBehavior(s, c.input)

			serv := services.AbstractService{AuthService: s}
			h := Handler{service: &serv}

			r := gin.New()
			r.POST("/sign-up", h.signUp)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(c.body))

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			} else {
				assert.Equal(t, rec.Body.String(), `{"id":1}`)
			}
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthService, input dto.SignInDTO)

	cases := []struct {
		name                string
		body                string
		input               dto.SignInDTO
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name: "OK",
			body: `{"username":"username","password":"password"}`,
			input: dto.SignInDTO{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignInDTO) {
				s.EXPECT().SignIn(input).Return(int64(1), nil)
				s.EXPECT().GenerateTokens(int64(1)).Return(gomock.Any().String(), gomock.Any().String(), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "Invalid body",
			body:                "invalid_body",
			input:               dto.SignInDTO{},
			mockBehavior:        func(s *mock_services.MockAuthService, input dto.SignInDTO) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "Invalid username or password",
			body: `{"username":"invalid_username","password":"invalid_password"}`,
			input: dto.SignInDTO{
				Username: "invalid_username",
				Password: "invalid_password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignInDTO) {
				s.EXPECT().SignIn(input).Return(int64(0), sql.ErrNoRows)
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
		{
			name: "Service failed",
			body: `{"username":"username","password":"password"}`,
			input: dto.SignInDTO{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignInDTO) {
				s.EXPECT().SignIn(input).Return(int64(0), errors.New("some error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
		{
			name: "Failed to generate tokens",
			body: `{"username":"username","password":"password"}`,
			input: dto.SignInDTO{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_services.MockAuthService, input dto.SignInDTO) {
				s.EXPECT().SignIn(input).Return(int64(1), nil)
				s.EXPECT().GenerateTokens(int64(1)).Return("", "", errors.New("some error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: true,
		},
	}

	cfg := &config.Config{
		Cookie: config.Cookie{
			Name:     "refresh-token",
			Age:      1000,
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_services.NewMockAuthService(ctrl)
			c.mockBehavior(s, c.input)

			serv := services.AbstractService{AuthService: s}
			h := NewHandler(&serv, cfg, new(memory.Cache))

			r := gin.New()
			r.POST("/sign-in", h.signIn)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(c.body))

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			} else {
				cookie := rec.Result().Cookies()
				assert.Len(t, cookie, 1)

				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["Bearer"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			}
		})
	}
}

func TestHandler_refresh(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthService, token string)

	cases := []struct {
		name                string
		cookie              http.Cookie
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name: "OK",
			cookie: http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				MaxAge:   1000,
				Path:     "/",
				Domain:   "/localhost",
				Secure:   false,
				HttpOnly: true,
			},
			mockBehavior: func(s *mock_services.MockAuthService, token string) {
				s.EXPECT().UpdateTokens(token).Return(gomock.Any().String(), gomock.Any().String(), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "Empty cookie",
			cookie:              http.Cookie{},
			mockBehavior:        func(s *mock_services.MockAuthService, token string) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "Invalid cookie",
			cookie: http.Cookie{
				Name:     "invalid-cookie",
				Value:    "token",
				MaxAge:   1000,
				Path:     "/",
				Domain:   "/localhost",
				Secure:   false,
				HttpOnly: true,
			},
			mockBehavior:        func(s *mock_services.MockAuthService, token string) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name: "Service failed",
			cookie: http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				MaxAge:   1000,
				Path:     "/",
				Domain:   "/localhost",
				Secure:   false,
				HttpOnly: true,
			},
			mockBehavior: func(s *mock_services.MockAuthService, token string) {
				s.EXPECT().UpdateTokens(token).Return("", "", errors.New("some error"))
			},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
	}

	cfg := &config.Config{
		Cookie: config.Cookie{
			Name:     "refresh-token",
			Age:      1000,
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_services.NewMockAuthService(ctrl)
			c.mockBehavior(s, c.cookie.Value)

			serv := services.AbstractService{AuthService: s}
			h := NewHandler(&serv, cfg, new(memory.Cache))

			r := gin.New()
			r.GET("/refresh", h.refresh)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/refresh", nil)
			req.AddCookie(&c.cookie)

			r.ServeHTTP(rec, req)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			} else {
				cookies := rec.Result().Cookies()
				assert.Len(t, cookies, 1)

				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["Bearer"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			}
		})
	}
}

func TestHandler_middlewareAuth(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthService, token string)

	cases := []struct {
		name                string
		header              string
		token               string
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedErrResponse bool
	}{
		{
			name:   "OK",
			header: "Authorization",
			token:  "Bearer token",
			mockBehavior: func(s *mock_services.MockAuthService, token string) {
				s.EXPECT().ParseToken(token).Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:                "Emty header",
			header:              "",
			token:               "Bearer token",
			mockBehavior:        func(s *mock_services.MockAuthService, token string) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:                "Invalid header",
			header:              "invalid_header",
			token:               "Bearer token",
			mockBehavior:        func(s *mock_services.MockAuthService, token string) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:                "Incorrect token",
			header:              "Authorization",
			token:               "invalid token",
			mockBehavior:        func(s *mock_services.MockAuthService, token string) {},
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: true,
		},
		{
			name:   "Invalid token",
			header: "Authorization",
			token:  "Bearer token",
			mockBehavior: func(s *mock_services.MockAuthService, token string) {
				s.EXPECT().ParseToken(token).Return(int64(0), errors.New("some error"))
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedErrResponse: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_services.NewMockAuthService(ctrl)
			c.mockBehavior(s, "token")

			serv := services.AbstractService{AuthService: s}
			h := Handler{service: &serv}

			r := gin.New()
			r.GET("/auth", h.middlewareAuth)

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			ctx.Request = httptest.NewRequest("GET", "/auth", nil)
			ctx.Request.Header.Set(c.header, c.token)

			h.middlewareAuth(ctx)

			assert.Equal(t, rec.Code, c.expectedStatus)
			if c.expectedErrResponse {
				var responseBody map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				value, ok := responseBody["message"]

				assert.True(t, ok)
				assert.NotEmpty(t, value)
			} else {
				id, _ := ctx.Get("user_id")
				assert.Equal(t, id, int64(1))
			}
		})
	}

}
