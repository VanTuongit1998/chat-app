package handler

import (
	"encoding/json"
	"net/http"

	"chat-app/internal/service"
	"chat-app/internal/utils"
)

type AuthHandler struct {
	authUsecase *service.AuthUsecase
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type userInfo struct {
	Username string `json:"username"`
}

func NewAuthHandler(authUsecase *service.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := h.authUsecase.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := h.authUsecase.GenerateToken(user)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "cannot generate token")
		return
	}

	utils.JSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) Users(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	users := h.authUsecase.ListUsers()
	result := make([]userInfo, 0, len(users))
	for _, user := range users {
		result = append(result, userInfo{Username: user.Username})
	}

	utils.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := h.authUsecase.Register(req.Username, req.Password)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]interface{}{"username": user.Username, "id": user.ID})
}
