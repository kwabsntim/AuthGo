package handlers

import (
	"AuthGo/models"
	"AuthGo/services"
	"encoding/json"
	"net/http"
)

// handlers/handlers.go
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}
	// Call service - let it handle all business logic
	user, err := services.RegisterUser(req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is nil
	if user == nil {
		http.Error(w, "User creation returned nil", http.StatusInternalServerError)
		return
	}

	// Create token
	token, err := services.CreateToken(user.ID.Hex())
	if err != nil {
		http.Error(w, "Token creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Success response
	response := models.AuthResponse{
		Message: "User created successfully",
		Token:   token,
		User: models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
