package handler

import (
	"net/http"
	"strconv"

	"chat-app/internal/middleware"
	"chat-app/internal/service"
	"chat-app/internal/utils"
	"chat-app/internal/websocket"
)

type ChatHandler struct {
	hub         *websocket.Hub
	chatUsecase *service.ChatUsecase
}

func NewChatHandler(hub *websocket.Hub, chatUsecase *service.ChatUsecase) *ChatHandler {
	return &ChatHandler{hub: hub, chatUsecase: chatUsecase}
}

func (h *ChatHandler) Websocket(w http.ResponseWriter, r *http.Request) {
	websocket.ServeWs(h.hub, w, r)
}

func (h *ChatHandler) Messages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := middleware.GetAuthenticatedUser(r)
	if user == nil || user.Username == "" {
		utils.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	contact := r.URL.Query().Get("with")
	if contact == "" {
		utils.Error(w, http.StatusBadRequest, "contact is required")
		return
	}

	limit := int64(100)
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil || parsedLimit <= 0 {
			utils.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = parsedLimit
	}

	messages, err := h.chatUsecase.Conversation(r.Context(), user.Username, contact, limit)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "cannot load messages")
		return
	}

	utils.JSON(w, http.StatusOK, messages)
}
