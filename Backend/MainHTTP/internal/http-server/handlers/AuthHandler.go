package handlers

import (
	"encoding/json"
	"fmt"
	"mainHTTP/internal/domain/models"
	"net/http"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error("Handlers.LoginUser at JSON: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.client.Login(r.Context(), user.Email, user.PassHash)
	if err != nil {
		h.logger.Error("Handlers.LoginUser: " + err.Error())
		http.Error(w, "Login error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error("Handlers.RegisterUser at JSON: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(user)

	token, err := h.client.Register(r.Context(), user.Email, string(user.PassHash))
	if err != nil {
		h.logger.Error("Handlers.RegisterUser: " + err.Error())
		http.Error(w, "Registration error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
