package services

import "AuthGo/models"

// RegisterInterface defines the contract for user registration operations
type RegisterInterface interface {
	// RegisterUser creates a new user with the provided credentials
	RegisterUser(email, username, password string) (*models.User, error)
}

// LoginInterface defines the contract for user authentication operations
type LoginInterface interface {
	// LoginUser authenticates a user with email and password
	LoginUser(email, password string) (*models.User, error)
}
type FetchUsersInterface interface {
	FetchAllUsers() ([]models.User, error)
}
