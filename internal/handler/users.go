package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// GetUserByID — обработчик для получения пользователя по его ID.
// Требует, чтобы ID был передан как параметр маршрута (`/v1.0/users/{id}`).
// Возвращает JSON с данными пользователя или ошибку:
//   - 404, если ID невалиден или пользователь не найден
//   - 400, если произошла ошибка сериализации
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.logger.Log.Warn("Невалидный ID пользователя", zap.String("id_param", idParam), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	readDTO, err := h.users.GetUserByID(r.Context(), id)
	if err != nil {
		h.logger.Log.Warn("Пользователь не найден", zap.Uint64("user_id", id), zap.Error(err))
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(readDTO)
	if err != nil {
		h.logger.Log.Error("Ошибка сериализации ответа", zap.Uint64("user_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Log.Info("Пользователь успешно получен", zap.Uint64("user_id", id))

	_, err = w.Write(resp)
	if err != nil {
		h.logger.Log.Error("Ошибка отправки ответа клиенту", zap.Uint64("user_id", id), zap.Error(err))
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}
