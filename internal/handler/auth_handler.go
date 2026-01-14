package handler

import (
	"Learning/internal/handler/requests"
	"Learning/internal/helper"
	"Learning/internal/model"
	"Learning/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	passwordHash, err := h.hashPassword(req.Password)
	if err != nil {
		return
	}
	u := model.User{
		Login:    req.Login,
		Name:     req.Name,
		Password: passwordHash,
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
		slog.Error("Not correct request body", "error", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		slog.Error("Failed to generate new tokens", "error", err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	h.setTokens(w, accessToken, refreshToken)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		slog.Error("Refresh token does not exists")
		http.Error(w, "no refresh token", http.StatusUnauthorized)
		return
	}
	accessToken, newRefreshToken, err := h.authService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		slog.Error("Failed to generate new tokens", "error", err)
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

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
	if err != nil {
		return
	}
}

func (h *AuthHandler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
