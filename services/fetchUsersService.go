package services

import (
	"AuthGo/models"
	"AuthGo/repository"
	"errors"
)

type FetchUserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewFetchUserService(userRepo repository.UserRepository) FetchUsersInterface {
	return &FetchUserServiceImpl{userRepo: userRepo}
}
func (s *FetchUserServiceImpl) FetchAllUsers() ([]models.User, error) {
	//fetch all users
	users, err := s.userRepo.FetchAllUsers()
	if err != nil {
		return nil, errors.New("could not fetch users")
	}
	return users, nil

}
