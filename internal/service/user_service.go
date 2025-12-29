package service

import (
	"Learning/internal/model"
	"Learning/internal/repository"
	"context"
	"errors"
)

var ErrUserAlreadyExists = errors.New("user with such login already exists")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	existed, err := s.repo.GetByLogin(ctx, user.Login)
	if err == nil && existed != nil {
		return ErrUserAlreadyExists
	}
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetById(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}
