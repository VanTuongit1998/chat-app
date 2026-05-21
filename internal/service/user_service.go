package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) ListUsers() []*model.User {
	return s.userRepo.FindAll()
}
