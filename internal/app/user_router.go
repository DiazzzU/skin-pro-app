package app

import (
	"skin-pro-app/internal/auth"
	"skin-pro-app/internal/config"
	"skin-pro-app/internal/handler"

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
