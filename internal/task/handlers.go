package task

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

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	task, err := h.service.FindOne(r.Context(), "_id", id)
	if err != nil {
		h.logger.Error("Task not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	ownerID, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	tasks, err := h.service.FindMany(r.Context(), "owner_id", ownerID)

	if err != nil {
		h.logger.Error("Tasks not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) GetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	if title == "" {
		err := fmt.Errorf("Title required")
		h.logger.Error("error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.FindOne(r.Context(), "title", title)
	if err != nil {
		h.logger.Error("Task not found", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	dto := &CreateTaskDTO{}

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		h.logger.Error("Failed to unmarshal req body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(r.Context(), dto)
	if err != nil {
		h.logger.Error("Failed to create task", "error", err)
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

func (h *handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	dto := &UpdateTaskDTO{}
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

	data, err := json.Marshal(fmt.Sprintf("Task with id = %s was updated", id.Hex()))
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err, status := h.getIDFromRequest(r)
	if err != nil {
		h.logger.Error("error", err)
		http.Error(w, err.Error(), status)
		return
	}

	_, err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete task", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fmt.Sprintf("Successfully deleted task"))
	if err != nil {
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
