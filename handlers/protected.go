package handlers

import (
	"AuthGo/middleware"
	"encoding/json"
	"net/http"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value(middleware.UserIDKey).(string)
	
	response := map[string]interface{}{
		"message": "Protected route accessed successfully",
		"user_id": userID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}