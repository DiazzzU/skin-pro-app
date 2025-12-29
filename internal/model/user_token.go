package model

import "time"

type UserToken struct {
	ID           int64     `json:"id"`
	RefreshToken string    `json:"refresh_token"`
	UserID       int64     `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}
