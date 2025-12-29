package app

import (
	"Learning/internal/handler"

	"github.com/go-chi/chi/v5"
)

func NewAuthRouter(userHandler *handler.AuthHandler) AuthRouter {
	r := chi.NewRouter()
	r.Post("/register", userHandler.Register) // POST /auth/register
	r.Post("/login", userHandler.Login)       // POST /auth/login
	r.Post("/refresh", userHandler.Refresh)   // POST /auth/refresh
	return r
}
