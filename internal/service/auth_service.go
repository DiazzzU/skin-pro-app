package service

import (
	"Learning/internal/auth"
	"Learning/internal/config"
	"Learning/internal/model"
	"Learning/internal/repository"
	"context"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepo      *repository.UserRepository
	userTokenRepo *repository.UserTokenRepository
	cfg           *config.GlobalConfig
}

func NewAuthService(u *repository.UserRepository, ut *repository.UserTokenRepository, cfg *config.GlobalConfig) *AuthService {
	return &AuthService{userRepo: u, userTokenRepo: ut, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (accessToken string, refreshToken string, err error) {
	slog.Info("Trying to login", "username", username)
	user, err := s.userRepo.GetByLogin(ctx, username)
	if err != nil || !s.checkPasswordHash(password, user.Password) {
		slog.Error("Invalid credentials", "username", username)
		return "", "", ErrInvalidCredentials
	}
	return s.generateToken(ctx, user.ID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, err error) {
	rt, err := s.userTokenRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		slog.Error("Failed to get refresh token", "error", err)
		return "", "", err
	}
	accessToken, newRefreshToken, err = s.generateToken(ctx, rt.UserID)
	_ = s.userTokenRepo.Revoke(ctx, refreshToken)
	return accessToken, newRefreshToken, err
}

func (s *AuthService) generateToken(ctx context.Context, userID int64) (accessToken string, refreshToken string, err error) {
	accessToken, err = auth.GenerateJWT(userID, s.cfg.JWTSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, expiresAt, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	ut := model.UserToken{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	if err := s.userTokenRepo.Create(ctx, &ut); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
