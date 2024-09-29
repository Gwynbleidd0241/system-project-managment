package handlers

import (
	authService "authService/internal/services/auth"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	router  *mux.Router
	service *authService.Auth
	RegisterHandler
	LoginHandler
}

type RegisterHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}

type LoginHandler interface {
	LoginUser(w http.ResponseWriter, r *http.Request)
}

func (h *Handler) RegistNewHandlers() {
	h.router.HandleFunc("/register", h.RegisterUser).Methods(http.MethodPost)
	h.router.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func NewHandler(router *mux.Router, auth *authService.Auth) *Handler {
	handler := &Handler{
		router:  router,
		service: auth,
	}
	handler.RegistNewHandlers()
	return handler
}
