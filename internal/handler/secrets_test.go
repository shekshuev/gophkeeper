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
	os.Setenv("ACCESS_TOKEN_SECRET", "test")
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	handler := NewHandler(nil, nil, secrets, &cfg)
	server := httptest.NewServer(handler.Router)
	defer server.Close()

	accessToken, _ := utils.CreateToken(cfg.AccessTokenSecret, "1", cfg.AccessTokenExpires)

	t.Run("Success", func(t *testing.T) {
		secrets.EXPECT().
			GetByID(gomock.Any(), gomock.Any()).
			Return(&models.ReadSecretDTO{
				ID:     1,
				UserID: 1,
				Title:  "My secret",
			}, nil)

		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessToken).
			Get(server.URL + "/v1.0/secrets/1")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("Invalid_ID", func(t *testing.T) {
		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessToken).
			Get(server.URL + "/v1.0/secrets/abc")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})

	t.Run("Secret_not_found", func(t *testing.T) {
		secrets.EXPECT().
			GetByID(gomock.Any(), gomock.Any()).
			Return(nil, assert.AnError)

		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessToken).
			Get(server.URL + "/v1.0/secrets/2")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})

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

	t.Run("Success_same_user_ID", func(t *testing.T) {
		dto := models.CreateSecretDTO{
			UserID: 1,
			Title:  "Secret A",
			Data:   models.SecretDataDTO{Text: ptr("top secret")},
		}

		secrets.EXPECT().
			Create(gomock.Any(), dto).
			Return(uint64(42), nil)

		body, _ := json.Marshal(dto)

		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessTokenForUser1).
			SetHeader("Content-Type", "application/json").
			SetBody(bytes.NewReader(body)).
			Post(server.URL + "/v1.0/secrets/")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
	})

	t.Run("Unauthorized_no_token", func(t *testing.T) {
		dto := models.CreateSecretDTO{
			UserID: 1,
			Title:  "Secret B",
			Data:   models.SecretDataDTO{Text: ptr("hacked")},
		}

		body, _ := json.Marshal(dto)

		resp, err := resty.New().R().
			SetHeader("Content-Type", "application/json").
			SetBody(bytes.NewReader(body)).
			Post(server.URL + "/v1.0/secrets/")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("Invalid_JSON", func(t *testing.T) {
		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessTokenForUser1).
			SetHeader("Content-Type", "application/json").
			SetBody(`{invalid json`).
			Post(server.URL + "/v1.0/secrets/")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
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

	accessTokenUser10, _ := utils.CreateToken(cfg.AccessTokenSecret, "10", cfg.AccessTokenExpires)
	accessTokenUser11, _ := utils.CreateToken(cfg.AccessTokenSecret, "11", cfg.AccessTokenExpires)

	t.Run("Success_same_user_ID", func(t *testing.T) {
		secrets.EXPECT().
			GetAllByUser(gomock.Any(), uint64(10)).
			Return([]models.ReadSecretDTO{
				{ID: 1, UserID: 10, Title: "First"},
				{ID: 2, UserID: 10, Title: "Second"},
			}, nil)

		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessTokenUser10).
			Get(server.URL + "/v1.0/secrets/user/10")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("Unauthorized_no_token", func(t *testing.T) {
		resp, err := resty.New().R().
			Get(server.URL + "/v1.0/secrets/user/10")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("Forbidden_token_subject_not_equal_user_id", func(t *testing.T) {
		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessTokenUser11).
			Get(server.URL + "/v1.0/secrets/user/10")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("Invalid_user_id_format", func(t *testing.T) {
		resp, err := resty.New().R().
			SetHeader("Authorization", "Bearer "+accessTokenUser10).
			Get(server.URL + "/v1.0/secrets/user/abc")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())
	})
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
