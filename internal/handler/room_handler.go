package handler

import (
	"net/http"

	"chat-app/internal/service"
	"chat-app/internal/utils"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) Rooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rooms, err := h.roomService.ListRooms(r.Context())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "cannot load rooms")
		return
	}

	utils.JSON(w, http.StatusOK, rooms)
}
