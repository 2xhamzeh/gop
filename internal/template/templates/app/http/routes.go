package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	userHandler *userHandler,
	noteHandler *noteHandler,
	authMiddleware func(http.Handler) http.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(notFound)
	r.MethodNotAllowed(methodNotAllowed)

	r.Use(
		requestID,
		loggerMiddleware,
		recovery,
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/user/register", userHandler.register)
			r.Post("/user/login", userHandler.login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/user", userHandler.getUser)
			r.Put("/user", userHandler.updateUser)
			r.Delete("/user", userHandler.deleteUser)

			r.Get("/notes", noteHandler.getUserNotes)
			r.Post("/notes", noteHandler.createNote)
			r.Put("/notes/{id}", noteHandler.updateNote)
			r.Delete("/notes/{id}", noteHandler.deleteNote)
		})
	})

	return r
}
