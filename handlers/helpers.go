package handlers

import (
	"AuthGo/models"
	"AuthGo/utils"
	"encoding/json"
	"net/http"
)

// DecodeJSONRequest decodes JSON request body into the provided destination
// Returns true if successful, false if error (and sends error response)
func DecodeJSONRequest(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

// CreateTokenForUser creates a JWT token for the given user
// Returns token string and error
func CreateTokenForUser(user *models.User) (string, error) {
	return utils.CreateToken(user.ID.Hex())
}

// SendAuthResponse sends a standardized authentication response
func SendAuthResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	response := models.JSONresponse{
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
