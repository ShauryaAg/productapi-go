package routes

import (
	"github.com/ShauryaAg/ProductAPI/handlers"
	"github.com/ShauryaAg/ProductAPI/middlewares"
	"github.com/go-chi/chi/v5"
)

func GetProductRoutes(r chi.Router) {
	r.Route("/product", func(r chi.Router) {
		r.With(middlewares.AuthMiddleware).Route("/", func(r chi.Router) {
			r.Route("/", func(r chi.Router) {
				r.Post("/", handlers.CreateProduct)
			})
		})
		r.Get("/", handlers.SearchProducts)
	})
}
