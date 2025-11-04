package handlers

import (
	"AuthGo/models"
	"AuthGo/utils"
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"regexp"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Client *mongo.Client
}

const (
	minPasswordLength = 8
	maxPasswordLength = 20
)

// creating the login handler
func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	//decoding the request body that is parsing json
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}
	//-----------Email actions------------------
	//checking for empty email
	if user.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}
	//checking the email against regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	//checking the password lenght
	passwordLength := utf8.RuneCountInString(user.Password)
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		http.Error(w, "Password length must be between 8 and 20 characters", http.StatusBadRequest)
		return
	}
	//parsing the email
	if _, err := mail.ParseAddress(user.Email); err != nil {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	//hashing the password
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.CreatedAt = time.Now()
	//-----------Database actions-------------------------------
	//inserting the user into the database
	collection := h.Client.Database("usersdb").Collection("profiles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//checking if the user already exists
	existing := collection.FindOne(ctx, bson.M{"email": user.Email})
	if existing.Err() == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	//creating user
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	// Use the result to send success response
	response := models.Message{Message: "User created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}
