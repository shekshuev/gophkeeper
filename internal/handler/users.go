package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetUserByID — обработчик для получения пользователя по его ID.
// Требует, чтобы ID был передан как параметр маршрута (`/v1.0/users/{id}`).
// Возвращает JSON с данными пользователя или ошибку:
//   - 404, если ID невалиден или пользователь не найден
//   - 400, если произошла ошибка сериализации
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}
	readDTO, err := h.users.GetUserByID(r.Context(), id)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}
	resp, err := json.Marshal(readDTO)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}
