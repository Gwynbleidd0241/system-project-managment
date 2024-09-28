package handlers

import (
	"authService/internal/domain/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user by providing email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} string "token"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Registration error"
// @Router /register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.RegisterNewUser(r.Context(), user.Email, string(user.PassHash))
	if err != nil {
		http.Error(w, "Registration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
