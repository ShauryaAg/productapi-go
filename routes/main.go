package routes

import "github.com/go-chi/chi/v5"

func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		GetUserRoutes(r)
	})

	return r
}
