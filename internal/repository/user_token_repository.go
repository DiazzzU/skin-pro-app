package repository

import (
	"Learning/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreateUserTokenQuery   = "INSERT INTO user_tokens (refresh_token, user_id, expires_at) VALUES ($1, $2, $3)"
	GetByRefreshTokenQuery = "SELECT id, refresh_token, user_id, expires_at FROM user_tokens WHERE refresh_token=$1 AND expires_at > now()"
	RevokeUserTokenQuery   = "DELETE FROM user_tokens WHERE refresh_token=$1"
)

type UserTokenRepository struct {
	db *pgxpool.Pool
}

func NewUserTokenRepository(db *pgxpool.Pool) *UserTokenRepository {
	return &UserTokenRepository{db: db}
}

func (r *UserTokenRepository) Create(ctx context.Context, token *model.UserToken) error {
	_, err := r.db.Exec(ctx, CreateUserTokenQuery, token.RefreshToken, token.UserID, token.ExpiresAt)
	return err
}

func (r *UserTokenRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.UserToken, error) {
	var userToken model.UserToken
	err := r.db.QueryRow(ctx,
		GetByRefreshTokenQuery, refreshToken,
	).Scan(&userToken.ID, &userToken.RefreshToken, &userToken.UserID, &userToken.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &userToken, err
}

func (r *UserTokenRepository) Revoke(ctx context.Context, refreshToken string) error {
	_, err := r.db.Exec(ctx, RevokeUserTokenQuery, refreshToken)
	return err
}
