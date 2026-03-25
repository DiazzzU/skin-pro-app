package handler

import (
	"encoding/json"
	"net/http"
	"skin-pro-app/internal/auth"
	"skin-pro-app/internal/config"
	"skin-pro-app/internal/handler/responses"
	"skin-pro-app/internal/service"
)

type UserHandler struct {
	svc *service.UserService
	cfg *config.GlobalConfig
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	u, err := h.svc.GetById(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	userInfo := responses.UserInfo{
		ID:    u.ID,
		Login: u.Login,
		Name:  u.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		return
	}
}
