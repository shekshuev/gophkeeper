package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/shekshuev/gophkeeper/internal/logger"
	"github.com/shekshuev/gophkeeper/internal/middleware"
	"github.com/shekshuev/gophkeeper/internal/service"
	"github.com/shekshuev/gophkeeper/internal/utils"
	"go.uber.org/zap"
)

// Handler — основной HTTP-обработчик приложения.
// Содержит зависимости: сервисы, валидатор, конфигурацию и chi.Router.
//
// Регистрирует маршруты:
//   - /v1.0/auth/login    — POST: логин пользователя
//   - /v1.0/auth/register — POST: регистрация пользователя
//   - /v1.0/users/{id}    — GET: получение пользователя по ID (требует JWT)
type Handler struct {
	users    service.UserService
	secrets  service.SecretService
	auth     service.AuthService
	Router   *chi.Mux
	validate *validator.Validate
	cfg      *config.Config
	logger   *logger.Logger
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var ErrValidationError = errors.New("validation error")
var ErrInvalidID = errors.New("invalid ID")
var ErrInvalidToken = errors.New("invalid token")
var ErrNotFound = errors.New("not found")

// NewHandler создаёт и настраивает HTTP-обработчик со всеми маршрутами и middleware.
// Использует:
//   - стандартные middleware chi (RequestID, Logger, Recoverer и др.)
//   - CORS (разрешает все источники)
//   - JWT-аутентификацию для защищённых маршрутов
func NewHandler(
	users service.UserService,
	auth service.AuthService,
	secrets service.SecretService,
	cfg *config.Config,
) *Handler {
	router := chi.NewRouter()
	validate := utils.NewValidator()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))
	router.Use(chiMiddleware.Recoverer)
	router.Use(cors.AllowAll().Handler)
	h := &Handler{users: users, auth: auth, secrets: secrets, Router: router, validate: validate, cfg: cfg, logger: logger.NewLogger()}

	h.Router.Route("/v1.0/users", func(r chi.Router) {
		r.Use(middleware.RequestAuth(cfg.AccessTokenSecret))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUserByID)
		})
	})

	h.Router.Route("/v1.0/secrets", func(r chi.Router) {
		r.With(middleware.RequestAuth(cfg.AccessTokenSecret)).Post("/", h.CreateSecret)
		r.With(middleware.RequestAuth(cfg.AccessTokenSecret)).Get("/{id:[0-9]+}", h.GetSecretByID)
		r.With(middleware.RequestAuth(cfg.AccessTokenSecret)).Delete("/{id:[0-9]+}", h.DeleteSecretByID)
		r.With(middleware.RequestAuthSameID(cfg.AccessTokenSecret)).Get("/user/{user_id:[0-9]+}", h.GetAllSecretsByUserID)
	})

	h.Router.Route("/v1.0/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})

	h.Router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			h.logger.Log.Error("failed to write health response", zap.Error(err))
		}
	})

	return h
}

// JSONError отправляет ответ с ошибкой в формате JSON:
//
//	{ "error": "<текст ошибки>" }
//
// Используется во всех обработчиках при ошибках на уровне валидации, авторизации и прочих.
func (h *Handler) JSONError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: errMessage})
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
}
