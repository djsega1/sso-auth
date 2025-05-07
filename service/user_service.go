package service

import (
	"fmt"

	"github.com/djsega1/sso-auth/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(username, password string) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %v", err)
	}

	return s.repo.CreateUser(username, string(hashedPassword))
}

func (s *UserService) AuthenticateUser(username, password string) (bool, error) {
	hashedPassword, err := s.repo.GetHashedPassword(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil // invalid password
	}

	return true, nil
}
