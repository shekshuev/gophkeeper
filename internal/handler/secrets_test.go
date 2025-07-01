package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetSecretByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secrets := mocks.NewMockSecretService(ctrl)
	cfg := config.GetConfig()
	handler := NewHandler(nil, nil, secrets, &cfg)
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	accessToken, _ := utils.CreateToken(cfg.AccessTokenSecret, "1", cfg.AccessTokenExpires)

	testCases := []struct {
		name          string
		secretID      string
		expectedCode  int
		responseDTO   *models.ReadSecretDTO
		serviceError  error
		serviceCalled bool
	}{
		// {
		// 	name:         "Success",
		// 	secretID:     "1",
		// 	expectedCode: http.StatusOK,
		// 	responseDTO: &models.ReadSecretDTO{
		// 		ID:     1,
		// 		UserID: 1,
		// 		Title:  "My secret",
		// 	},
		// 	serviceCalled: true,
		// },
		// {
		// 	name:          "Invalid ID",
		// 	secretID:      "abc",
		// 	expectedCode:  http.StatusNotFound,
		// 	serviceCalled: false,
		// },
		{
			name:          "Secret not found",
			secretID:      "2",
			expectedCode:  http.StatusNotFound,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				secrets.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(tc.responseDTO, tc.serviceError)
			}

			resp, err := resty.New().R().
				SetHeader("Authorization", "Bearer "+accessToken).
				Get(server.URL + "/v1.0/secrets/" + tc.secretID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_CreateSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secrets := mocks.NewMockSecretService(ctrl)
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()

	handler := NewHandler(nil, nil, secrets, &cfg)
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	accessTokenForUser1, _ := utils.CreateToken(cfg.AccessTokenSecret, "1", cfg.AccessTokenExpires)

	testCases := []struct {
		name          string
		accessToken   string
		payload       models.CreateSecretDTO
		expectedCode  int
		serviceCalled bool
	}{
		{
			name:        "Success (same user ID)",
			accessToken: accessTokenForUser1,
			payload: models.CreateSecretDTO{
				UserID: 1,
				Title:  "Secret A",
				Data:   models.SecretDataDTO{Text: ptr("top secret")},
			},
			expectedCode:  http.StatusCreated,
			serviceCalled: true,
		},
		{
			name:        "Unauthorized (no token)",
			accessToken: "",
			payload: models.CreateSecretDTO{
				UserID: 1,
				Title:  "Secret B",
				Data:   models.SecretDataDTO{Text: ptr("hacked")},
			},
			expectedCode:  http.StatusUnauthorized,
			serviceCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				secrets.EXPECT().
					Create(gomock.Any(), tc.payload).
					Return(uint64(42), nil)
			}

			body, _ := json.Marshal(tc.payload)

			resp, err := resty.New().R().
				SetHeader("Authorization", "Bearer "+tc.accessToken).
				SetHeader("Content-Type", "application/json").
				SetBody(bytes.NewReader(body)).
				Post(server.URL + "/v1.0/secrets/")

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_GetAllSecretsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secrets := mocks.NewMockSecretService(ctrl)
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()

	handler := NewHandler(nil, nil, secrets, &cfg)
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	// accessTokenUser10, _ := utils.CreateToken(cfg.AccessTokenSecret, "10", cfg.AccessTokenExpires)
	accessTokenUser11, _ := utils.CreateToken(cfg.AccessTokenSecret, "11", cfg.AccessTokenExpires)

	testCases := []struct {
		name          string
		requestUserID string
		accessToken   string
		expectedCode  int
		serviceCalled bool
	}{
		// {
		// 	name:          "Success (same user ID)",
		// 	requestUserID: "10",
		// 	accessToken:   accessTokenUser10,
		// 	expectedCode:  http.StatusOK,
		// 	serviceCalled: true,
		// },
		// {
		// 	name:          "Unauthorized (no token)",
		// 	requestUserID: "10",
		// 	accessToken:   "",
		// 	expectedCode:  http.StatusUnauthorized,
		// 	serviceCalled: false,
		// },
		{
			name:          "Forbidden (token subject ≠ user_id)",
			requestUserID: "10",
			accessToken:   accessTokenUser11,
			expectedCode:  http.StatusUnauthorized,
			serviceCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				secrets.EXPECT().
					GetAllByUser(gomock.Any(), uint64(10)).
					Return([]models.ReadSecretDTO{
						{ID: 1, UserID: 10, Title: "First"},
						{ID: 2, UserID: 10, Title: "Second"},
					}, nil)
			}

			req := resty.New().R()
			if tc.accessToken != "" {
				req.SetHeader("Authorization", "Bearer "+tc.accessToken)
			}

			resp, err := req.
				Get(server.URL + "/v1.0/secrets/user/" + tc.requestUserID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
		})
	}
}

func TestHandler_DeleteSecretByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secrets := mocks.NewMockSecretService(ctrl)
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()

	accessToken, _ := utils.CreateToken(cfg.AccessTokenSecret, "77", cfg.AccessTokenExpires)

	handler := NewHandler(nil, nil, secrets, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	secrets.EXPECT().
		DeleteByID(gomock.Any(), uint64(77)).
		Return(nil)

	resp, err := resty.New().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Delete(httpSrv.URL + "/v1.0/secrets/77")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode())
}

func ptr(s string) *string {
	return &s
}
