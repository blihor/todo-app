package user

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	user, err := h.service.FindOne(r.Context(), "_id", id)
	if err != nil {
		h.logger.Error("User not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	dto := &CreateUserDTO{}

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		h.logger.Error("Failed to unmarshal req body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(result.InsertedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	_, err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fmt.Sprintf("Successfully deleted user"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	dto := &UpdateUserDTO{}
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		h.logger.Error("Failed to unmarshal req body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.Update(r.Context(), id, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fmt.Sprintf("User with id = %s was updated", id.Hex()))
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) getIDFromRequest(r *http.Request) (bson.ObjectID, error, int) {
	id := r.PathValue("id")
	if id == "" {
		err := fmt.Errorf("ID hasn't been specified")
		return bson.NilObjectID, err, http.StatusBadRequest
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, fmt.Errorf("Failed to convert id to ObjectID type %w", err), http.StatusInternalServerError
	}

	return objID, nil, http.StatusOK
}
