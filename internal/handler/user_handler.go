package handler

import (
	"net/http"

	"chat-app/internal/service"
	"chat-app/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	users := h.userService.ListUsers()
	result := make([]userInfo, 0, len(users))
	for _, user := range users {
		result = append(result, userInfo{Username: user.Username})
	}

	utils.JSON(w, http.StatusOK, result)
}
