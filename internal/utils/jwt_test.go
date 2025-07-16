package utils

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	secret := "testsecret"
	userId := "12345"
	exp := time.Hour

	token, err := CreateToken(secret, userId, exp)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.Nil(t, err)
	assert.NotNil(t, parsedToken)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "Gophkeeper", claims["iss"])
	assert.Equal(t, userId, claims["sub"])
}

func TestGetToken(t *testing.T) {
	secret := "testsecret"
	userId := "12345"
	exp := time.Second * 1

	token, err := CreateToken(secret, userId, exp)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := GetToken(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userId, claims.Subject)

	time.Sleep(2 * exp)

	claims, err = GetToken(token, secret)
	assert.ErrorIs(t, err, ErrTokenExpired)
	assert.Nil(t, claims)
}

func TestGetToken_InvalidSignature(t *testing.T) {
	token, err := CreateToken("correct-secret", "user", time.Minute)
	assert.NoError(t, err)

	claims, err := GetToken(token, "wrong-secret")
	assert.ErrorIs(t, err, ErrInvalidSignature)
	assert.Nil(t, claims)
}

func TestGetToken_InvalidFormat(t *testing.T) {
	claims, err := GetToken("not.a.jwt", "secret")
	assert.ErrorIs(t, err, ErrTokenInvalid)
	assert.Nil(t, claims)
}

func TestGetRawAccessToken(t *testing.T) {
	req := &http.Request{Header: http.Header{}}

	token, err := GetRawAccessToken(req)
	assert.ErrorIs(t, err, ErrTokenInvalid)
	assert.Empty(t, token)

	req.Header.Set("Authorization", "Token somevalue")
	_, err = GetRawAccessToken(req)
	assert.ErrorIs(t, err, ErrTokenInvalid)

	req.Header.Set("Authorization", "Bearer sometoken")
	token, err = GetRawAccessToken(req)
	assert.NoError(t, err)
	assert.Equal(t, "sometoken", token)
}

func TestGetRawRefreshToken(t *testing.T) {
	req := &http.Request{Header: http.Header{}}

	token, err := GetRawRefreshToken(req)
	assert.Error(t, err)
	assert.Empty(t, token)

	req.AddCookie(&http.Cookie{
		Name:  RefreshTokenCookieName,
		Value: "refresh123",
	})
	token, err = GetRawRefreshToken(req)
	assert.NoError(t, err)
	assert.Equal(t, "refresh123", token)
}

func TestContextClaims(t *testing.T) {
	claims := jwt.RegisteredClaims{Subject: "user123"}

	ctx := context.Background()
	ctxWithClaims := PutClaimsToContext(ctx, claims)

	extracted, ok := GetClaimsFromContext(ctxWithClaims)
	assert.True(t, ok)
	assert.Equal(t, claims.Subject, extracted.Subject)

	_, ok = GetClaimsFromContext(ctx)
	assert.False(t, ok)
}
