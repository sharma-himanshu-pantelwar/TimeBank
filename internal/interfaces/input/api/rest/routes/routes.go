package routes

import (
	"net/http"
	userhandler "timebank/internal/interfaces/input/api/rest/handler"
	"timebank/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *userhandler.UserHandler) http.Handler {
	router := chi.NewRouter()
	router.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
	})
	router.Route("/v1/user", func(r chi.Router) {
		r.Use(middleware.Authenticate)

		r.Get("/profile", userHandler.Profile)
		r.Post("/logout", userHandler.Logout)
	})
	router.Route("/v1/skills", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Post("/add", userHandler.AddSkills)
		r.Get("/find/{skill}", userHandler.FindSkilledPerson)
		r.Patch("/rename/{skillId}", userHandler.RenameSkill)
		r.Delete("/delete/{skillId}", userHandler.DeleteSkill)
	})

	return router
}
