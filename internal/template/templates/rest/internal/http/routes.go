package http

import (
	"github.com/go-chi/chi/v5"
)

// NewRouter creates a new router, registers all routes and middlewares and returns the router.
// It uses chi as the underlying router.
func NewRouter(
	userHandler *userHandler,
	middlewares *middlewares,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(
		middlewares.RequestID,
		middlewares.Logger,
		middlewares.Recovery,
	)

	r.NotFound(middlewares.NotFound)

	r.MethodNotAllowed(middlewares.MethodNotAllowed)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/user/register", userHandler.register)
		r.Post("/user/login", userHandler.login)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth)

			r.Get("/user", userHandler.getUser)
			r.Patch("/user", userHandler.updateUser)
			r.Delete("/user", userHandler.deleteUser)
		})
	})

	return r
}
