package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRequestAuth_Success(t *testing.T) {
	secret := "secret"
	token, err := utils.CreateToken(secret, "42", time.Minute)
	assert.NoError(t, err)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		claims, ok := utils.GetClaimsFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "42", claims.Subject)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	mw := RequestAuth(secret)
	mw(handler).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, called)
}

func TestRequestAuth_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	resp := httptest.NewRecorder()

	mw := RequestAuth("secret")
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRequestAuthSameID_Success(t *testing.T) {
	secret := "secret"
	userID := "123"
	token, err := utils.CreateToken(secret, userID, time.Minute)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", userID)

	req := httptest.NewRequest("GET", "/users/"+userID, nil).WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	called := false
	mw := RequestAuthSameID(secret)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, called)
}

func TestRequestAuthSameID_Unauthorized(t *testing.T) {
	secret := "secret"
	token, err := utils.CreateToken(secret, "123", time.Minute)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", "999")

	req := httptest.NewRequest("GET", "/users/999", nil).WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	mw := RequestAuthSameID(secret)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
