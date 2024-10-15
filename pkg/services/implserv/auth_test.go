package implserv

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/core/entity"
	"github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/config"
	mock_repositories "github.com/DmytroBeliasnyk/crud_app_rest_api/pkg/repositories/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockAuthRepository(ctrl)

	serv := AuthServiceImpl{repo: repo}

	input := dto.SignUpDTO{
		Name:     "Name",
		Email:    "aaa@bbb.ccc",
		Username: "username",
		Password: "password",
	}

	expectedUser := entity.FromSignUpDTO(input, serv.HashPassword(input.Password))

	repo.EXPECT().SignUp(expectedUser).Return(int64(1), nil)

	got, err := serv.SignUp(input)

	assert.NoError(t, err)
	assert.Equal(t, got, int64(1))
}

func TestAuthService_SignIn(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockAuthRepository, username, passwordHash string)

	cases := []struct {
		name         string
		input        dto.SignInDTO
		mockBehavior mockBehavior
		expected     int64
		expectedErr  bool
	}{
		{
			name: "OK",
			input: dto.SignInDTO{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(s *mock_repositories.MockAuthRepository, username, passwordHash string) {
				s.EXPECT().SignIn(username, passwordHash).Return(int64(1), nil)
			},
			expected: 1,
		},
		{
			name: "Invalid data",
			input: dto.SignInDTO{
				Username: "invalid_username",
				Password: "invalid_password",
			},
			mockBehavior: func(s *mock_repositories.MockAuthRepository, username, passwordHash string) {
				s.EXPECT().SignIn(username, passwordHash).Return(int64(0), sql.ErrNoRows)
			},
			expectedErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockAuthRepository(ctrl)
			serv := AuthServiceImpl{
				repo: repo,
			}
			c.mockBehavior(repo, c.input.Username, serv.HashPassword(c.input.Password))

			got, err := serv.SignIn(c.input)
			if c.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, c.expected)
			}
		})
	}
}

func TestAuthService_GenerateTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockAuthRepository(ctrl)

	repo.
		EXPECT().
		CreateRefreshToken(int64(1), gomock.Any(), gomock.Any()).
		Return(nil)

	cfg := &config.Config{
		Auth: config.Auth{
			Signature: "signature",
			JWT:       time.Hour,
			Refresh:   time.Hour * 2,
		},
	}

	jt, rt, err := NewAuthService(repo, cfg).GenerateTokens(1)

	assert.NoError(t, err)
	assert.NotEmpty(t, jt)
	assert.NotEmpty(t, rt)
}

func TestAuthService_UpdateTokens(t *testing.T) {
	type mockBehavior func(s *mock_repositories.MockAuthRepository, token string)

	cases := []struct {
		name         string
		input        string
		mockBehavior mockBehavior
		expectedErr  bool
	}{
		{
			name:  "OK",
			input: "token",
			mockBehavior: func(s *mock_repositories.MockAuthRepository, token string) {
				s.EXPECT().FindRefreshToken(token).Return(int64(1),
					time.Date(2025, time.January, 0, 0, 0, 0, 0, time.Local), nil)
				s.EXPECT().DeleteRefreshToken(token).Return(nil)
				s.EXPECT().CreateRefreshToken(int64(1), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "Failed to generate",
			input: "token",
			mockBehavior: func(s *mock_repositories.MockAuthRepository, token string) {
				s.EXPECT().FindRefreshToken(token).Return(int64(1),
					time.Date(2025, time.January, 0, 0, 0, 0, 0, time.Local), nil)
				s.EXPECT().DeleteRefreshToken(token).Return(nil)
				s.EXPECT().CreateRefreshToken(int64(1), gomock.Any(), gomock.Any()).
					Return(errors.New("some error"))
			},
			expectedErr: true,
		},
		{
			name:  "Invalid token",
			input: "invalid_token",
			mockBehavior: func(s *mock_repositories.MockAuthRepository, token string) {
				s.EXPECT().FindRefreshToken(token).Return(int64(0), time.Time{}, sql.ErrNoRows)
			},
			expectedErr: true,
		},
		{
			name:  "Expired token",
			input: "expired_token",
			mockBehavior: func(s *mock_repositories.MockAuthRepository, token string) {
				s.EXPECT().FindRefreshToken(token).Return(int64(1), time.Now().AddDate(0, 0, -1), nil)
				s.EXPECT().DeleteRefreshToken(token).Return(nil)
			},
			expectedErr: true,
		},
	}

	cfg := &config.Config{
		Auth: config.Auth{
			Signature: "signature",
			JWT:       time.Hour,
			Refresh:   time.Hour * 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repositories.NewMockAuthRepository(ctrl)
			c.mockBehavior(repo, c.input)

			jt, rt, err := NewAuthService(repo, cfg).UpdateTokens(c.input)
			if c.expectedErr {
				assert.Error(t, err)
				assert.Empty(t, jt)
				assert.Empty(t, rt)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, jt)
				assert.NotEmpty(t, rt)
			}
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	mockGenerateToken := func(id string, signingMethod jwt.SigningMethod,
		signature string, issuedAt, expiresAt time.Time) string {
		jwtt := jwt.NewWithClaims(signingMethod, &jwt.RegisteredClaims{
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		})

		jt, _ := jwtt.SignedString([]byte(signature))
		return jt
	}

	cases := []struct {
		name        string
		token       string
		expected    int64
		expectedErr bool
	}{
		{
			name: "Valid token",
			token: mockGenerateToken("1", jwt.SigningMethodHS256,
				"signature", time.Now(), time.Now().Add(time.Hour)),
			expected: 1,
		},
		{
			name: "Invalid signing method",
			token: mockGenerateToken("1", jwt.SigningMethodES256,
				"signature", time.Now(), time.Now().Add(time.Hour)),
			expectedErr: true,
		},
		{
			name: "Invalid signature",
			token: mockGenerateToken("1", jwt.SigningMethodHS256,
				"invalid_signature", time.Now(), time.Now().Add(time.Hour)),
			expectedErr: true,
		},
		{
			name: "Invalid subject",
			token: mockGenerateToken("invalid_id", jwt.SigningMethodHS256,
				"signature", time.Now(), time.Now().Add(time.Hour)),
			expectedErr: true,
		},
		{
			name: "Expired token",
			token: mockGenerateToken("1", jwt.SigningMethodHS256,
				"signature", time.Now().AddDate(0, 0, -2), time.Now().AddDate(0, 0, -1)),
			expectedErr: true,
		},
	}

	cfg := &config.Config{
		Auth: config.Auth{
			Signature: "signature",
		},
	}

	serv := NewAuthService(new(mock_repositories.MockAuthRepository), cfg)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := serv.ParseToken(c.token)
			if c.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, id, c.expected)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id, c.expected)
			}
		})
	}
}
