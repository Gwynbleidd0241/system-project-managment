package handlers

import (
	"encoding/json"
	"net/http"
	"taskService/internal/domain/models"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, "Fail to Create Task", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.URL.Query().Get("taskName")
	err := h.service.DeleteTask(r.Context(), taskName)
	if err != nil {
		http.Error(w, "Fail to Delete Task", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, "Fail to Update Task", http.StatusInternalServerError)
		return
	}
}
