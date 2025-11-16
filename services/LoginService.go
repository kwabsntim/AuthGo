package services

import (
	"AuthGo/models"
	"AuthGo/repository"
	"AuthGo/utils"
	"AuthGo/validation"
	"errors"
)

// this stores everything from the user repository
// all the functions from user_repo.go are called here
type LoginServiceImpl struct {
	userRepo repository.UserRepository
}

// this puts userRepo into the LoginServiceImpl struct and this
// is the link to call the methods from the repository
func NewLoginService(userRepo repository.UserRepository) LoginInterface {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}

	return &LoginServiceImpl{userRepo: userRepo}
}

func (s *LoginServiceImpl) LoginUser(email, password string) (*models.User, error) {
	//validating the input
	err := validation.ValidateEmail(email)
	if err != nil {
		return nil, err
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	//getting the user email from repo
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	//comparing the password
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
