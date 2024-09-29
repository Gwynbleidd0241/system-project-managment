package handlers

import (
	"authService/internal/domain/models"
	"encoding/json"
	"net/http"
)

// LoginUser godoc
// @Summary Login user
// @Description Log in a user by email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User login details"
// @Success 200 {string} string "token"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Login error"
// @Router /login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), user.Email, string(user.PassHash))
	if err != nil {
		http.Error(w, "Login error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
