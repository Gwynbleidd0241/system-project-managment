package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"taskService/internal/service"
)

type Handler struct {
	router  *mux.Router
	service *service.TaskService
	TaskHandler
}

type TaskHandler interface {
	CreateTask(writer http.ResponseWriter, request *http.Request)
	UpdateTask(writer http.ResponseWriter, request *http.Request)
	DeleteTask(writer http.ResponseWriter, request *http.Request)
	GetTasks(writer http.ResponseWriter, request *http.Request)
}

func (h *Handler) RegistNewHandlers() {
	h.router.HandleFunc("/task", h.CreateTask).Methods(http.MethodPost)
	h.router.HandleFunc("/task", h.UpdateTask).Methods(http.MethodPut)
	//h.router.HandleFunc("/task", h.GetTasks).Methods(http.MethodGet)
	h.router.HandleFunc("/task", h.DeleteTask).Methods(http.MethodDelete)
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func NewHandler(router *mux.Router, auth *service.TaskService) *Handler {
	handler := &Handler{
		router:  router,
		service: auth,
	}
	handler.RegistNewHandlers()
	return handler
}
