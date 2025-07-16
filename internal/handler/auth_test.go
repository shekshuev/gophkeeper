package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/mocks"
	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	auth := mocks.NewMockAuthService(ctrl)
	cfg := config.GetConfig()
	handler := NewHandler(nil, auth, nil, &cfg)

	t.Run("Success login", func(t *testing.T) {
		dto := models.LoginUserDTO{UserName: "test_user", Password: "test123!"}
		auth.EXPECT().Login(gomock.Any(), dto).Return(&models.ReadTokenDTO{AccessToken: "access", RefreshToken: "refresh"}, nil)
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Wrong password", func(t *testing.T) {
		dto := models.LoginUserDTO{UserName: "test_user", Password: "WrongPassword123!"}
		auth.EXPECT().Login(gomock.Any(), dto).Return(nil, service.ErrWrongPassword)
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		dto := models.LoginUserDTO{UserName: "ghost", Password: "test123!"}
		auth.EXPECT().Login(gomock.Any(), dto).Return(nil, service.ErrUserNotFound)
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Validation error - short user name", func(t *testing.T) {
		dto := models.LoginUserDTO{UserName: "usr", Password: "test123!"}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Validation error - bad password", func(t *testing.T) {
		dto := models.LoginUserDTO{UserName: "test_user", Password: "123"}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Body read error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", io.NopCloser(&brokenReader{}))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/login", strings.NewReader(`{"user_name": "test_user", "password":}`))
		rr := httptest.NewRecorder()
		handler.Login(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

}

func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	auth := mocks.NewMockAuthService(ctrl)
	cfg := config.GetConfig()
	handler := NewHandler(nil, auth, nil, &cfg)

	t.Run("Success register", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "test_user", Password: "test123!", PasswordConfirm: "test123!", FirstName: "John", LastName: "Doe",
		}
		auth.EXPECT().Register(gomock.Any(), dto).Return(&models.ReadTokenDTO{AccessToken: "access", RefreshToken: "refresh"}, nil)
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("Password mismatch", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "test_user", Password: "12345678", PasswordConfirm: "1234", FirstName: "John", LastName: "Doe",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Validation error - short user name", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "usr", Password: "test123!", PasswordConfirm: "test123!", FirstName: "John", LastName: "Doe",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Validation error - weak password", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "test_user", Password: "123", PasswordConfirm: "123", FirstName: "John", LastName: "Doe",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Validation error - empty FirstName", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "test_user", Password: "test123!", PasswordConfirm: "test123!", FirstName: "", LastName: "Doe",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Validation error - empty LastName", func(t *testing.T) {
		dto := models.RegisterUserDTO{
			UserName: "test_user", Password: "test123!", PasswordConfirm: "test123!", FirstName: "John", LastName: "",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("Body read error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", io.NopCloser(&brokenReader{}))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1.0/auth/register", strings.NewReader(`{"user_name": "test_user", "password":}`))
		rr := httptest.NewRecorder()
		handler.Register(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

}

type brokenReader struct {
	called bool
}

func (r *brokenReader) Read(p []byte) (int, error) {
	if !r.called {
		r.called = true
		return 0, io.ErrUnexpectedEOF
	}
	return 0, io.EOF
}

func (r *brokenReader) Close() error {
	return nil
}

type FakeReadTokenDTO struct{}

func (FakeReadTokenDTO) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("mock marshal error")
}
