package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type handler struct {
	service *service
	logger  *slog.Logger
}

func NewHandler(service *service, logger *slog.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	dto := &UserLoginDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to unmarshal req body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err, status := h.service.Login(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to sign a jwt token", "error", err)
		http.Error(w, err.Error(), status)
		return
	}

	data, err := json.Marshal(token)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	dto := &UserRegisterDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Error("Failed to unmarshal req body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err, status := h.service.Register(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		http.Error(w, err.Error(), status)
		return
	}

	data, err := json.Marshal(id)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
