package service

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthServiceImpl(t *testing.T) {
	cfg := config.GetConfig()
	repo := mocks.NewMockUserRepository(gomock.NewController(t))
	svc := NewAuthServiceImpl(repo, &cfg)
	assert.NotNil(t, svc)
}

func TestAuthServiceImpl_Login(t *testing.T) {
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("REFRESH_TOKEN_SECRET", "test")
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	authService := &AuthServiceImpl{repo: repo, cfg: &cfg, logger: logger.NewLogger()}
	ctx := context.Background()

	testCases := []struct {
		name     string
		dto      models.LoginUserDTO
		hasError bool
		mockSet  func()
	}{
		{
			name: "Success",
			dto: models.LoginUserDTO{
				UserName: "testuser",
				Password: "password123",
			},
			hasError: false,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName(ctx, "testuser").Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					PasswordHash: utils.HashPassword("password123"),
				}, nil)
			},
		},
		{
			name: "User not found",
			dto: models.LoginUserDTO{
				UserName: "nonexistentuser",
				Password: "password123",
			},
			hasError: true,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName(ctx, "nonexistentuser").Return(nil, sql.ErrNoRows)
			},
		},
		{
			name: "Wrong password",
			dto: models.LoginUserDTO{
				UserName: "testuser",
				Password: "wrongpassword",
			},
			hasError: true,
			mockSet: func() {
				repo.EXPECT().GetUserByUserName(ctx, "testuser").Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					PasswordHash: utils.HashPassword("password123"),
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSet()
			_, err := authService.Login(ctx, tc.dto)
			if tc.hasError {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestAuthServiceImpl_Register(t *testing.T) {
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("REFRESH_TOKEN_SECRET", "test")
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	authService := &AuthServiceImpl{repo: repo, cfg: &cfg, logger: logger.NewLogger()}
	ctx := context.Background()

	fixedPasswordHash := "$2a$10$CmIxNqxCFrgFoji4qyka0.UvTV4wG54LN5UJjV7mfH6q0caiNGUvK"

	testCases := []struct {
		name        string
		dto         models.RegisterUserDTO
		expectedErr bool
		mockSet     func()
	}{
		{
			name: "Success",
			dto: models.RegisterUserDTO{
				UserName:        "testuser",
				Password:        "password123",
				PasswordConfirm: "password123",
				FirstName:       "Test",
				LastName:        "User",
			},
			expectedErr: false,
			mockSet: func() {
				repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(&models.ReadAuthUserDataDTO{
					ID:           1,
					UserName:     "testuser",
					PasswordHash: fixedPasswordHash,
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSet()
			_, err := authService.Register(ctx, tc.dto)
			if tc.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestAuthServiceImpl_Register_CreateUserError(t *testing.T) {
	cfg := config.GetConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	authService := NewAuthServiceImpl(repo, &cfg)
	ctx := context.Background()

	dto := models.RegisterUserDTO{
		UserName:        "testuser",
		Password:        "password123",
		PasswordConfirm: "password123",
		FirstName:       "Test",
		LastName:        "User",
	}

	repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, assert.AnError)

	_, err := authService.Register(ctx, dto)
	assert.Error(t, err)
}

func TestAuthServiceImpl_generateTokenPair_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)

	cfg := &config.Config{
		AccessTokenSecret:   "",
		RefreshTokenSecret:  "",
		AccessTokenExpires:  0,
		RefreshTokenExpires: 0,
	}

	service := NewAuthServiceImpl(repo, cfg)
	user := models.ReadAuthUserDataDTO{
		ID:       1,
		UserName: "brokenuser",
	}

	token, err := service.generateTokenPair(user)
	assert.Error(t, err)
	assert.Nil(t, token)
}
