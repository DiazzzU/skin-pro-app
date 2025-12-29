package handler

import (
	"Learning/internal/handler/requests"
	"Learning/internal/helper"
	"Learning/internal/model"
	"Learning/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, err := helper.DecodeJSONBody[requests.RegisterRequest](w, r)
	if err != nil {
		return
	}
	u := model.User{
		Login:    req.Login,
		Name:     req.Name,
		Password: req.Password,
	}
	if err := h.userService.Create(r.Context(), &u); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			http.Error(w, "User with such login already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := helper.DecodeJSONBody[requests.LoginRequest](w, r)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	h.setTokens(w, accessToken, refreshToken)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "no refresh token", http.StatusUnauthorized)
		return
	}
	accessToken, newRefreshToken, err := h.authService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh", http.StatusUnauthorized)
		return
	}
	h.setTokens(w, accessToken, newRefreshToken)
}

func (h *AuthHandler) setTokens(w http.ResponseWriter, accessToken string, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/auth/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	err := json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
	if err != nil {
		return
	}
}
