package services

import "AuthGo/models"

type UserService interface {
	RegisterUser(email, username, password string) (*models.User, error)
}