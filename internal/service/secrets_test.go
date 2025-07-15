package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSecretServiceImpl_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	now := time.Now()

	testCases := []struct {
		name      string
		id        uint64
		setupMock func()
		expected  *models.ReadSecretDTO
		hasError  bool
	}{
		{
			name: "Success",
			id:   1,
			setupMock: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), uint64(1)).
					Return(&models.ReadSecretDTO{
						ID:        1,
						UserID:    10,
						Title:     "my note",
						Data:      models.SecretDataDTO{Text: ptr("hello")},
						CreatedAt: now,
						UpdatedAt: now,
					}, nil)
			},
			expected: &models.ReadSecretDTO{
				ID:     1,
				UserID: 10,
				Title:  "my note",
				Data:   models.SecretDataDTO{Text: ptr("hello")},
			},
			hasError: false,
		},
		{
			name: "Not found",
			id:   404,
			setupMock: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), uint64(404)).
					Return(nil, errors.New("not found"))
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			result, err := service.GetByID(context.Background(), tc.id)

			if tc.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Title, result.Title)
				assert.Equal(t, *tc.expected.Data.Text, *result.Data.Text)
			}
		})
	}
}

func TestSecretServiceImpl_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	input := models.CreateSecretDTO{
		UserID: 10,
		Title:  "Secret",
		Data:   models.SecretDataDTO{Text: ptr("text")},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			Create(gomock.Any(), input).
			Return(uint64(123), nil)

		id, err := service.Create(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), id)
	})
}

func TestSecretServiceImpl_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	input := models.CreateSecretDTO{
		UserID: 10,
		Title:  "FailSecret",
		Data:   models.SecretDataDTO{Text: ptr("fail")},
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), input).
		Return(uint64(0), errors.New("insert error"))

	id, err := service.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, uint64(0), id)
}

func TestSecretServiceImpl_GetAllByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetAllByUser(gomock.Any(), uint64(10)).
			Return([]models.ReadSecretDTO{
				{ID: 1, Title: "A"},
				{ID: 2, Title: "B"},
			}, nil)

		result, err := service.GetAllByUser(context.Background(), 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "A", result[0].Title)
	})
}

func TestSecretServiceImpl_GetAllByUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	mockRepo.EXPECT().
		GetAllByUser(gomock.Any(), uint64(10)).
		Return(nil, errors.New("db error"))

	result, err := service.GetAllByUser(context.Background(), 10)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSecretServiceImpl_DeleteByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			DeleteByID(gomock.Any(), uint64(77)).
			Return(nil)

		err := service.DeleteByID(context.Background(), 77)
		assert.NoError(t, err)
	})
}

func TestSecretServiceImpl_DeleteByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSecretRepository(ctrl)
	service := NewSecretServiceImpl(mockRepo)

	mockRepo.EXPECT().
		DeleteByID(gomock.Any(), uint64(77)).
		Return(errors.New("delete error"))

	err := service.DeleteByID(context.Background(), 77)
	assert.Error(t, err)
}

func ptr(s string) *string {
	return &s
}
