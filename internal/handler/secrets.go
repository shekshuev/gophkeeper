package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/shekshuev/gophkeeper/internal/models"
	"github.com/shekshuev/gophkeeper/internal/utils"
)

// GetSecretByID — обработчик для получения секрета по его ID.
// Возвращает JSON с данными секрета или ошибку:
//   - 404, если ID невалиден или секрет не найден
//   - 400, если произошла ошибка сериализации
func (h *Handler) GetSecretByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Log.Warn("Невалидный ID секрета", zap.String("id", idStr), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	secret, err := h.secrets.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Log.Warn("Ошибка получения секрета", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	if secret == nil {
		h.logger.Log.Warn("Секрет не найден", zap.Uint64("secret_id", id))
		h.JSONError(w, http.StatusNotFound, ErrNotFound.Error())
		return
	}

	resp, err := json.Marshal(secret)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации секрета", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Log.Info("Секрет успешно получен", zap.Uint64("secret_id", id))
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка отправки секрета клиенту", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// GetAllSecretsByUserID — обработчик для получения всех секретов пользователя по его ID.
// Возвращает JSON с массивом секретов или ошибку:
//   - 404, если ID невалиден или ошибка при получении данных
//   - 400, если произошла ошибка сериализации
func (h *Handler) GetAllSecretsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.logger.Log.Warn("Невалидный user_id", zap.String("user_id", userIDStr), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	secrets, err := h.secrets.GetAllByUser(r.Context(), userID)
	if err != nil {
		h.logger.Log.Warn("Ошибка при получении секретов", zap.Uint64("user_id", userID), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(secrets)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации списка секретов", zap.Uint64("user_id", userID), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Log.Info("Секреты успешно получены", zap.Uint64("user_id", userID))
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка отправки ответа клиенту", zap.Uint64("user_id", userID), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// CreateSecret — обработчик создания нового секрета.
// Принимает JSON с полями title и data в теле запроса.
// Требует авторизации (по токену). Возвращает ID созданного секрета.
//
// Возвращает:
//   - 201 Created — если успешно
//   - 400 Bad Request — если JSON невалиден или ошибка сериализации
//   - 401 Unauthorized — если токен невалиден или не содержит userID
//   - 500 Internal Server Error — если ошибка на уровне сервиса
func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Log.Error("Не удалось прочитать тело запроса", zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, "cannot read body")
		return
	}
	defer r.Body.Close()

	var dto models.CreateSecretDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		h.logger.Log.Warn("Невалидный JSON при создании секрета", zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		h.logger.Log.Warn("Токен отсутствует или невалиден при создании секрета")
		h.JSONError(w, http.StatusUnauthorized, "missing or invalid token")
		return
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		h.logger.Log.Warn("Некорректный user ID в токене", zap.String("sub", claims.Subject), zap.Error(err))
		h.JSONError(w, http.StatusUnauthorized, "invalid user ID in token")
		return
	}

	dto.UserID = userID
	id, err := h.secrets.Create(r.Context(), dto)
	if err != nil {
		h.logger.Log.Error("Ошибка при создании секрета", zap.Uint64("user_id", userID), zap.Error(err))
		h.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Log.Info("Секрет успешно создан", zap.Uint64("secret_id", id), zap.Uint64("user_id", userID))

	resp := map[string]uint64{"id": id}
	encoded, err := json.Marshal(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации ID созданного секрета", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(encoded)
	if err != nil {
		h.logger.Log.Error("Ошибка отправки ответа при создании секрета", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// DeleteSecretByID — обработчик удаления секрета по ID.
// Возвращает:
//   - 204 No Content — если удаление прошло успешно
//   - 404 Not Found — если ID невалиден
//   - 500 Internal Server Error — если ошибка на уровне сервиса
func (h *Handler) DeleteSecretByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Log.Warn("Невалидный ID секрета при удалении", zap.String("id", idStr), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	err = h.secrets.DeleteByID(r.Context(), id)
	if err != nil {
		h.logger.Log.Error("Ошибка при удалении секрета", zap.Uint64("secret_id", id), zap.Error(err))
		h.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logger.Log.Info("Секрет успешно удалён", zap.Uint64("secret_id", id))
	w.WriteHeader(http.StatusNoContent)
}
