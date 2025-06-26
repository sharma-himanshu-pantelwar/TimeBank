package routes

import (
	"net/http"
	userhandler "timebank/internal/interfaces/input/api/rest/handler"
	"timebank/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *userhandler.UserHandler) http.Handler {
	router := chi.NewRouter()
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
	})
	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Get("/profile", userHandler.Profile)
		r.Post("/logout", userHandler.Logout)

	})

	return router
}
