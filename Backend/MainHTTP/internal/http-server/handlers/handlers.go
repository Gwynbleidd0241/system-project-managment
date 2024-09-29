package handlers

import (
	log "github.com/go-ozzo/ozzo-log"
	"github.com/gorilla/mux"
	"mainHTTP/internal/clients"
	"mainHTTP/internal/http-server/middleware/corse"
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
	h.router.Handle("/api/task", corse.New()(http.HandlerFunc(h.CreateTask))).Methods(http.MethodPost)
	h.router.Handle("/api/task", corse.New()(http.HandlerFunc(h.UpdateTask))).Methods(http.MethodPut)
	h.router.Handle("/api/task", corse.New()(http.HandlerFunc(h.GetTasks))).Methods(http.MethodGet)
	h.router.Handle("/api/task", corse.New()(http.HandlerFunc(h.DeleteTask))).Methods(http.MethodDelete)
	h.router.HandleFunc("/api/task", h.HandleOptions).Methods(http.MethodOptions)
	//h.router.HandleFunc("/api/task", h.CreateTask).Methods(http.MethodPost)
	//h.router.HandleFunc("/api/task", h.UpdateTask).Methods(http.MethodPut)
	//h.router.HandleFunc("/api/task", h.GetTasks).Methods(http.MethodGet)
	//h.router.HandleFunc("/api/task", h.DeleteTask).Methods(http.MethodDelete)
}

func (h *Handler) RegistAuthHandlers() {
	h.router.Handle("/api/login", corse.New()(http.HandlerFunc(h.LoginUser))).Methods(http.MethodPost, http.MethodOptions)
	h.router.Handle("/api/register", corse.New()(http.HandlerFunc(h.RegisterUser))).Methods(http.MethodPost, http.MethodOptions)
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func (h *Handler) HandleOptions(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
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
