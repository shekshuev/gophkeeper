package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceImpl_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	cfg := config.GetConfig()

	service := NewUserServiceImpl(mockRepo, &cfg)

	testCases := []struct {
		name      string
		id        uint64
		setupMock func()
		expected  *models.ReadUserDTO
		hasError  bool
	}{
		{
			name: "Success",
			id:   1,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&models.ReadUserDTO{
						ID:        1,
						UserName:  "john",
						FirstName: "John",
						LastName:  "Doe",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
			},
			expected: &models.ReadUserDTO{
				ID:        1,
				UserName:  "john",
				FirstName: "John",
				LastName:  "Doe",
			},
			hasError: false,
		},
		{
			name: "User not found",
			id:   999,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), uint64(999)).
					Return(nil, ErrUserNotFound)
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			tc.setupMock()

			result, err := service.GetUserByID(ctx, tc.id)

			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.UserName, result.UserName)
				assert.Equal(t, tc.expected.FirstName, result.FirstName)
				assert.Equal(t, tc.expected.LastName, result.LastName)
			}
		})
	}
}
