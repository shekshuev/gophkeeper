package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/gophkeeper/internal/utils"
)

// RequestAuth — middleware, проверяющий наличие и валидность access-токена в заголовке Authorization.
// Если токен валиден, добавляет claims в context.Context и передаёт управление следующему обработчику.
func RequestAuth(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.GetRawAccessToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, err := utils.GetToken(tokenString, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := utils.PutClaimsToContext(r.Context(), *claims)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestAuthSameID — middleware, проверяющий валидность токена и соответствие subject токена и ID в URL.
// Используется, когда доступ к ресурсу должен быть ограничен только его владельцем.
func RequestAuthSameID(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.GetRawAccessToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, err := utils.GetToken(tokenString, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			strId := chi.URLParam(r, "user_id")
			_, err = strconv.Atoi(strId)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			if strId != claims.Subject {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
