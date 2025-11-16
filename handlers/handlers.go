package handlers

import (
	"AuthGo/models"
	"AuthGo/services"
	"AuthGo/utils"
	"encoding/json"
	"log"
	"net/http"
)

type SignUpHandler struct {
	registerService services.RegisterInterface
}
type LoginHandler struct {
	loginService services.LoginInterface
}

func NewSignUpHandler(registerService services.RegisterInterface) *SignUpHandler {
	return &SignUpHandler{registerService: registerService}
}
func NewLoginHandler(loginService services.LoginInterface) *LoginHandler {
	return &LoginHandler{loginService: loginService}
}

func (h *SignUpHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.User

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call service - let it handle all business logic

	user, err := h.registerService.RegisterUser(req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Create token

	token, err := utils.CreateToken(user.ID.Hex())
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Login attempt for email: %s", req.Email)
	user, err := h.loginService.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := utils.CreateToken(user.ID.Hex())
	if err != nil {
		http.Error(w, "Token creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := models.AuthResponse{
		Message: "Login successful",
		Token:   token,
		User: models.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	responseBytes, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
	log.Printf("Login successful for user: %s", user.Email)
}