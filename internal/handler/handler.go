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
	"github.com/shekshuev/gophkeeper/internal/middleware"
	"github.com/shekshuev/gophkeeper/internal/service"
	"github.com/shekshuev/gophkeeper/internal/utils"
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
	auth     service.AuthService
	Router   *chi.Mux
	validate *validator.Validate
	cfg      *config.Config
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var ErrValidationError = errors.New("validation error")
var ErrInvalidID = errors.New("invalid ID")
var ErrInvalidToken = errors.New("invalid token")

// NewHandler создаёт и настраивает HTTP-обработчик со всеми маршрутами и middleware.
// Использует:
//   - стандартные middleware chi (RequestID, Logger, Recoverer и др.)
//   - CORS (разрешает все источники)
//   - JWT-аутентификацию для защищённых маршрутов /v1.0/users/{id}
func NewHandler(
	users service.UserService,
	auth service.AuthService,
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
	h := &Handler{users: users, auth: auth, Router: router, validate: validate, cfg: cfg}

	h.Router.Route("/v1.0/users", func(r chi.Router) {
		r.Use(middleware.RequestAuth(cfg.AccessTokenSecret))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetUserByID)
		})
	})

	h.Router.Route("/v1.0/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
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
