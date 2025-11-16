package services

import (
	"AuthGo/models"
	"AuthGo/repository"
	"AuthGo/utils"
	"AuthGo/validation"
	"fmt"
	"time"
)

// this stores all functions from the user repository
type RegisterServiceImpl struct {
	userRepo repository.UserRepository
}

func NewRegisterService(userRepo repository.UserRepository) RegisterInterface {
	return &RegisterServiceImpl{userRepo: userRepo}
}

func (s *RegisterServiceImpl) RegisterUser(email, username, password string) (*models.User, error) {
	//validating the input calling the validation functions
	start := time.Now()

	validationStart := time.Now()
	err := validation.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	err = validation.ValidatePassword(password)
	if err != nil {
		return nil, err
	}
	err = validation.ValidateUsername(username)
	if err != nil {
		return nil, err
	}
	fmt.Printf("    ⏱️ Validation: %v\n", time.Since(validationStart))
	//hashing password calling the utils function
	hashStart := time.Now()
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	fmt.Printf("    ⏱️ Password Hash: %v\n", time.Since(hashStart))
	//creating user body calling the models struct

	createStart := time.Now()
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Username: username,
	}
	fmt.Printf("    ⏱️ User Object: %v\n", time.Since(createStart))
	//saving the user by calling the respository function
	dbStart := time.Now()
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	fmt.Printf("    ⏱️ Database Save: %v\n", time.Since(dbStart))

	fmt.Printf("    ⏱️ RegisterUser Total: %v\n", time.Since(start))
	return user, nil
}
