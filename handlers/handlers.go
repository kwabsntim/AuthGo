package handlers

import (
	"AuthGo/models"
	"AuthGo/services"
	"encoding/json"
	"net/http"
)

// creating the login handler
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//decoding the request body
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//calling the service function to create the user
	user, err := services.RegisterUser(req.Email, req.Username, req.Password)
	//getting a token for the created user
	token, err := services.CreateToken(user.ID.Hex())
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	response := models.AuthResponse{
		Message: "User created successfully",
		Token:   token,
		User: models.UserResponse{
			ID:       user.ID.Hex(), // Convert ObjectID to string for JSON
			Username: user.Username,
			Email:    user.Email,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}
