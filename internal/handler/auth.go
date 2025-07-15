package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/shekshuev/gophkeeper/internal/models"
)

// Login — обработчик входа пользователя.
// Принимает JSON с полями user_name и password в теле запроса.
// Валидирует входные данные, вызывает auth-сервис и возвращает пару токенов.
//
// Возвращает:
//   - 200 OK — если авторизация прошла успешно
//   - 401 Unauthorized — если пароль неверен, пользователь не найден или ошибка парсинга/вызова сервиса
//   - 422 Unprocessable Entity — если входные данные не прошли валидацию
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.LoginUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err = json.Unmarshal(body, &loginDTO); err != nil {
		h.logger.Log.Error("Ошибка парсинга JSON", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.validate.Struct(loginDTO)
	if err != nil {
		h.logger.Log.Warn("Ошибка валидации входных данных", zap.Error(err))
		h.JSONError(w, http.StatusUnprocessableEntity, ErrValidationError.Error())
		return
	}
	tokensDTO, err := h.auth.Login(r.Context(), loginDTO)
	if err != nil {
		h.logger.Log.Warn("Ошибка входа пользователя", zap.String("user_name", loginDTO.UserName), zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp, err := json.Marshal(tokensDTO)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации токенов", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	h.logger.Log.Info("Пользователь успешно вошёл", zap.String("user_name", loginDTO.UserName))
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка при отправке ответа", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
	}
}

// Register — обработчик регистрации нового пользователя.
// Принимает JSON с user_name, password, password_confirm, first_name и last_name.
// Валидирует входные данные, вызывает auth-сервис и возвращает пару токенов.
//
// Возвращает:
//   - 201 Created — если регистрация прошла успешно
//   - 401 Unauthorized — если ошибка при создании пользователя или сериализации
//   - 422 Unprocessable Entity — если входные данные не прошли валидацию (в том числе проверку на совпадение паролей)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO models.RegisterUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Log.Error("Ошибка чтения тела запроса", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err = json.Unmarshal(body, &registerDTO); err != nil {
		h.logger.Log.Error("Ошибка парсинга JSON", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.validate.Struct(registerDTO)
	if err != nil {
		h.logger.Log.Warn("Ошибка валидации данных при регистрации", zap.Error(err))
		h.JSONError(w, http.StatusUnprocessableEntity, ErrValidationError.Error())
		return
	}
	tokensDTO, err := h.auth.Register(r.Context(), registerDTO)
	if err != nil {
		h.logger.Log.Error("Ошибка регистрации пользователя", zap.String("user_name", registerDTO.UserName), zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	resp, err := json.Marshal(tokensDTO)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации токенов", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	h.logger.Log.Info("Пользователь успешно зарегистрирован", zap.String("user_name", registerDTO.UserName))
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка при отправке ответа", zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, err.Error())
	}
}
