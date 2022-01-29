package routes

import (
	"github.com/ShauryaAg/ProductAPI/handlers"
	"github.com/ShauryaAg/ProductAPI/middlewares"
	"github.com/go-chi/chi/v5"
)

func GetReviewRoutes(r chi.Router) {
	r.Route("/review", func(r chi.Router) {
		r.With(middlewares.AuthMiddleware).Route("/{productId}", func(r chi.Router) {
			r.Post("/", handlers.CreateReview)
		})
	})
}
