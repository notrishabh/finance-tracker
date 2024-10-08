package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/notrishabh/finance-tracker/pkg/models"
	"github.com/notrishabh/finance-tracker/pkg/services"
	"github.com/notrishabh/finance-tracker/pkg/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if user.Password == "" || user.Username == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	if user.Password == "" || user.Username == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	loggedInUser, err := services.LoginUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(loggedInUser.Username, loggedInUser.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User logged in successfully",
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Only for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	json.NewEncoder(w).Encode(response)
}
