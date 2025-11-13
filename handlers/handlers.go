package handlers

import (
	"AuthGo/models"
	"AuthGo/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SignUpHandler struct {
	userService services.UserService
}

func NewSignUpHandler(userService services.UserService) *SignUpHandler {
	return &SignUpHandler{userService: userService}
}

func (h *SignUpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()

	var req models.User
	decodeStart := time.Now()
	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("⏱️ JSON Decode: %v\n", time.Since(decodeStart))
	if err != nil {
		fmt.Printf("❌ JSON Decode Error: %v\n", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("📝 Request data: email=%s, username=%s\n", req.Email, req.Username)
	// Call service - let it handle all business logic
	registerStart := time.Now()
	user, err := h.userService.RegisterUser(req.Email, req.Username, req.Password)
	fmt.Printf("⏱️ Register User: %v\n", time.Since(registerStart))
	if err != nil {
		fmt.Printf("❌ Registration Error: %v\n", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Create token
	tokenStart := time.Now()
	token, err := services.CreateToken(user.ID.Hex())
	fmt.Printf("⏱️ Create Token: %v\n", time.Since(tokenStart))
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
	responseStart := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	fmt.Printf("⏱️ Response: %v\n", time.Since(responseStart))
	fmt.Printf("🔥 Total Request: %v\n\n", time.Since(start))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
