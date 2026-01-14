package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserRouter http.Handler
type AuthRouter http.Handler

func NewRouter(userRouter UserRouter, authRouter AuthRouter) *chi.Mux {
	r := chi.NewRouter()
	r.Use(LoggingMiddleware)
	r.Mount("/users", userRouter)
	r.Mount("/auth", authRouter)
	return r
}
