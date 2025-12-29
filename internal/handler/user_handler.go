package handler

import (
	"Learning/internal/auth"
	"Learning/internal/config"
	"Learning/internal/service"
	"encoding/json"
	"net/http"
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		return
	}
}
