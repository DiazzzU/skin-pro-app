//go:build wireinject

package di

import (
	"skin-pro-app/internal/app"
	"skin-pro-app/internal/config"
	"skin-pro-app/internal/handler"
	"skin-pro-app/internal/repository"
	"skin-pro-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitUserRouter(pool *pgxpool.Pool, cfg *config.GlobalConfig) app.UserRouter {
	wire.Build(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,
		app.NewUserRouter,
	)
	return nil
}

func InitAuthRouter(pool *pgxpool.Pool, cfg *config.GlobalConfig) app.AuthRouter {
	wire.Build(
		repository.NewUserRepository,
		repository.NewUserTokenRepository,
		service.NewUserService,
		service.NewAuthService,
		handler.NewAuthHandler,
		app.NewAuthRouter,
	)
	return nil
}

func InitRouter(pool *pgxpool.Pool, cfg *config.GlobalConfig) *chi.Mux {
	wire.Build(
		InitUserRouter,
		InitAuthRouter,
		app.NewRouter,
	)
	return &chi.Mux{}
}
