package app

import (
	"Learning/internal/auth"
	"Learning/internal/config"
	"Learning/internal/handler"

	"github.com/go-chi/chi/v5"
)

func NewUserRouter(userHandler *handler.UserHandler, cfg *config.GlobalConfig) UserRouter {
	r := chi.NewRouter()

	r.Group(func(protected chi.Router) {
		protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
		protected.Get("/info", userHandler.GetUserInfo) // GET /users/{id}
	})
	return r
}
