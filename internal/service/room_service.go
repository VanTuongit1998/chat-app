package service

import (
	"context"

	"chat-app/internal/model"
	"chat-app/internal/repository"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
}

func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) ListRooms(ctx context.Context) ([]*model.Room, error) {
	return s.roomRepo.FindAll(ctx)
}
