package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/gophkeeper/internal/models"
)

// GetSecretByID — возвращает секрет по ID.
func (h *Handler) GetSecretByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	secret, err := h.secrets.GetByID(r.Context(), id)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	if secret == nil {
		h.JSONError(w, http.StatusNotFound, ErrNotFound.Error())
		return
	}

	resp, err := json.Marshal(secret)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// GetAllSecretsByUserID — возвращает все секреты пользователя по его userID.
func (h *Handler) GetAllSecretsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	secrets, err := h.secrets.GetAllByUser(r.Context(), userID)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(secrets)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// CreateSecret — создаёт новый секрет из тела запроса.
func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, "cannot read body")
		return
	}
	defer r.Body.Close()

	var dto models.CreateSecretDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		h.JSONError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	id, err := h.secrets.Create(r.Context(), dto)
	if err != nil {
		h.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]uint64{"id": id}
	encoded, err := json.Marshal(resp)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(encoded)
	if err != nil {
		h.JSONError(w, http.StatusBadRequest, err.Error())
	}
}

// DeleteSecretByID — удаляет секрет по ID.
func (h *Handler) DeleteSecretByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.JSONError(w, http.StatusNotFound, ErrInvalidID.Error())
		return
	}

	err = h.secrets.DeleteByID(r.Context(), id)
	if err != nil {
		h.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
