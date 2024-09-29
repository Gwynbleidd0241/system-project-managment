package handlers

import (
	"encoding/json"
	"fmt"
	authv1 "github.com/Dragzet/gRPCProtosv2/gen/go/task"
	"mainHTTP/internal/domain/models"
	"net/http"
)

func (h *Handler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task models.Task
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		h.logger.Error("Handlers.CreateTask at JSON: " + err.Error())
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcTask := &authv1.CreateTaskRequest{
		Email:       task.UserEmail,
		Name:        task.Name,
		Description: task.Description,
		Token:       task.Token,
	}

	id, err := h.client.CreateTask(request.Context(), grpcTask)
	if err != nil {
		h.logger.Error("Handlers.CreateTask: " + err.Error())
		http.Error(writer, "Create task error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"id": id})
}

func (h *Handler) DeleteTask(writer http.ResponseWriter, request *http.Request) {
	taskID := request.URL.Query().Get("id")

	grpcTask := &authv1.DeleteTaskRequest{
		Id: taskID,
	}

	status, err := h.client.DeleteTask(request.Context(), grpcTask)
	if err != nil {
		h.logger.Error("Handlers.DeleteTask: " + err.Error())
		http.Error(writer, "Delete task error", http.StatusInternalServerError)
		return
	}

	if status {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"status": "not found"})

	}
}

func (h *Handler) UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task models.Task
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		h.logger.Error("Handlers.UpdateTask at JSON: " + err.Error())
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	grpcTask := &authv1.UpdateTaskRequest{
		Email:       task.UserEmail,
		Name:        task.Name,
		Description: task.Description,
		Done:        task.IsDone,
		Token:       task.Token,
	}

	status, err := h.client.UpdateTask(request.Context(), grpcTask)
	if err != nil {
		h.logger.Error("Handlers.UpdateTask: " + err.Error())
		http.Error(writer, "Update task error", http.StatusInternalServerError)
		return
	}

	if status {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(map[string]string{"status": "not found"})

	}
}

func (h *Handler) GetTasks(writer http.ResponseWriter, request *http.Request) {
	email := request.URL.Query().Get("email")
	fmt.Println(email)
	grpcTasks, err := h.client.GetTasks(request.Context(), email)
	if err != nil {
		h.logger.Error("Handlers.GetTasks: " + err.Error())
		http.Error(writer, "Get tasks error", http.StatusInternalServerError)
		return
	}

	tasks := make([]models.Task, 0)
	for _, grpcTask := range grpcTasks {
		task := models.Task{
			ID:          grpcTask.Id,
			UserEmail:   grpcTask.Email,
			Name:        grpcTask.Name,
			Description: grpcTask.Description,
			IsDone:      grpcTask.Done,
		}
		tasks = append(tasks, task)
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(tasks)
}
