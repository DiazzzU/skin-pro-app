package repository

import (
	"Learning/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreateUserQuery  = `INSERT INTO users (login, name, password) VALUES ($1, $2, $3) RETURNING id`
	GetUserByIDQuery = `SELECT id, login, name, password FROM users WHERE id=$1`
	GetUserByLogin   = `SELECT id, login, name, password FROM users WHERE login=$1`
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	err := r.db.QueryRow(ctx, CreateUserQuery, u.Login, u.Name, u.Password).Scan(&u.ID)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx, GetUserByIDQuery, id).Scan(&u.ID, &u.Login, &u.Name, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx, GetUserByLogin, login).Scan(&u.ID, &u.Login, &u.Name, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}
