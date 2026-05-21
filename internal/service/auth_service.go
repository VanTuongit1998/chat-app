package service

import (
	"errors"

	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/internal/utils"
)

type AuthUsecase struct {
	userRepo   *repository.UserRepository
	jwtService *utils.JwtService
}

func NewAuthUsecase(userRepo *repository.UserRepository, jwtService *utils.JwtService) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtService: jwtService}
}

func (u *AuthUsecase) Authenticate(username, password string) (*model.User, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("username or password invalid")
	}
	return user, nil
}

func (u *AuthUsecase) ListUsers() []*model.User {
	return u.userRepo.FindAll()
}

func (u *AuthUsecase) GenerateToken(user *model.User) (string, error) {
	return u.jwtService.GenerateToken(user)
}

// Register creates a new user if username not taken
func (u *AuthUsecase) Register(username, password string) (*model.User, error) {
	// basic validation
	if username == "" || password == "" {
		return nil, errors.New("username and password required")
	}
	// ensure unique
	if _, err := u.userRepo.FindByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{Username: username, Password: password}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
