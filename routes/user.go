package routes

import (
	"github.com/ShauryaAg/ProductAPI/handlers"
	"github.com/ShauryaAg/ProductAPI/middlewares"

	"github.com/go-chi/chi/v5"
)

func GetUserRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handlers.Login)
		r.Post("/register", handlers.Register)

		r.With(middlewares.AuthMiddleware).Route("/user", func(r chi.Router) {
			r.Get("/", handlers.GetUser)
		})
	})

}
