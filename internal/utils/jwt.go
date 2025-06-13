package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// ContextKey используется для ключей в context.Context.
type ContextKey string

const (
	// AccessTokenCookieName — имя куки, в которой может храниться access-токен.
	AccessTokenCookieName = "X-Access-Token"

	// RefreshTokenCookieName — имя куки для refresh-токена.
	RefreshTokenCookieName = "X-Refresh-Token"

	// ContextClaimsKey — ключ для хранения JWT claims в context.Context.
	ContextClaimsKey = ContextKey("user-claims")
)

var (
	// ErrInvalidSignature возвращается, если подпись токена некорректна.
	ErrInvalidSignature = fmt.Errorf("token signature is invalid")

	// ErrTokenExpired возвращается, если токен истёк.
	ErrTokenExpired = fmt.Errorf("token is expired")

	// ErrTokenInvalid возвращается, если токен невалиден по другим причинам.
	ErrTokenInvalid = fmt.Errorf("token is invalid")
)

// GetRawAccessToken извлекает access-токен из заголовка Authorization.
// Ожидается формат: "Authorization: Bearer <token>".
func GetRawAccessToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrTokenInvalid
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// GetRawRefreshToken извлекает refresh-токен из cookie.
func GetRawRefreshToken(req *http.Request) (string, error) {
	cookie, err := req.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// CreateToken создаёт JWT-токен с указанным userId, сроком жизни и секретом.
// Возвращает подписанную строку токена.
func CreateToken(secret, userId string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "GopherTalk",
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	})
	return token.SignedString([]byte(secret))
}

// GetToken парсит и валидирует JWT по строке и секрету.
// Возвращает claims или соответствующую ошибку.
func GetToken(tokenString, secret string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, ErrInvalidSignature
			}
		}
		return nil, ErrTokenInvalid
	}

	if !token.Valid {
		return nil, ErrInvalidSignature
	}

	return claims, nil
}

// GetClaimsFromContext извлекает JWT claims из context.Context.
func GetClaimsFromContext(ctx context.Context) (jwt.RegisteredClaims, bool) {
	claims, ok := ctx.Value(ContextClaimsKey).(jwt.RegisteredClaims)
	return claims, ok
}

// PutClaimsToContext сохраняет JWT claims в context.Context.
func PutClaimsToContext(ctx context.Context, claims jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, ContextClaimsKey, claims)
}
