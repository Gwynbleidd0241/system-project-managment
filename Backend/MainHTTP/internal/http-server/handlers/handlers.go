package handlers

import (
	log "github.com/go-ozzo/ozzo-log"
	"github.com/gorilla/mux"
	"mainHTTP/internal/clients"
	"net/http"
)

type Handler struct {
	logger *log.Logger
	router *mux.Router
	client *clients.GRPCClient
	TaskHandler
	AuthHandler
}

type TaskHandler interface {
	CreateTask(writer http.ResponseWriter, request *http.Request)
	UpdateTask(writer http.ResponseWriter, request *http.Request)
	DeleteTask(writer http.ResponseWriter, request *http.Request)
	GetTasks(writer http.ResponseWriter, request *http.Request)
}

type AuthHandler interface {
	RegisterUser(writer http.ResponseWriter, request *http.Request)
	LoginUser(writer http.ResponseWriter, request *http.Request)
}

func (h *Handler) RegistNewTaskHandlers() {
	//h.router.HandleFunc("/task", h.CreateTask).Methods(http.MethodPost)
	//h.router.HandleFunc("/task", h.UpdateTask).Methods(http.MethodPut)
	//h.router.HandleFunc("/task", h.GetTasks).Methods(http.MethodGet)
	//h.router.HandleFunc("/task", h.DeleteTask).Methods(http.MethodDelete)
}

func (h *Handler) RegistAuthHandlers() {
	h.router.HandleFunc("/api/login", h.LoginUser).Methods(http.MethodPost)
	h.router.HandleFunc("/api/register", h.RegisterUser).Methods(http.MethodPost)
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func NewHandler(logger *log.Logger, router *mux.Router, client *clients.GRPCClient) *Handler {
	handler := &Handler{
		logger: logger,
		router: router,
		client: client,
	}
	handler.RegistAuthHandlers()
	handler.RegistNewTaskHandlers()
	return handler
}
