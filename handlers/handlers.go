package handlers

import (
	"AuthGo/models"
	"AuthGo/services"
	"AuthGo/utils"
	"net/http"
)

// ---------------------------Structs----------------------------
type SignUpHandler struct {
	registerService services.RegisterInterface
}
type LoginHandler struct {
	loginService services.LoginInterface
}
type FetchUsersHandler struct {
	fetchUsersService services.FetchUsersInterface
}

func NewSignUpHandler(registerService services.RegisterInterface) *SignUpHandler {
	return &SignUpHandler{registerService: registerService}
}
func NewLoginHandler(loginService services.LoginInterface) *LoginHandler {
	return &LoginHandler{loginService: loginService}
}
func NewFetchHandler(fetchUsersService services.FetchUsersInterface) *FetchUsersHandler {
	return &FetchUsersHandler{fetchUsersService: fetchUsersService}
}

//---------------------------Handler functions----------------------------

func (h *SignUpHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method

	var req models.User
	if !DecodeJSONRequest(w, r, &req) {
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
	data := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	// Success response
	SendAuthResponse(w, "User created successfully", data, http.StatusCreated)

}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req models.User
	//decoding the json request
	if !DecodeJSONRequest(w, r, &req) {
		return
	}

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
	data := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	// Success response
	SendAuthResponse(w, "Login successful", data, http.StatusCreated)
}
func (h *FetchUsersHandler) FetchAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.fetchUsersService.FetchAllUsers()
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"users": users,
	}
	SendAuthResponse(w, "Fetch successful", data, http.StatusCreated)
}
