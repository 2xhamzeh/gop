package http

import (
	"net/http"

	"example.com/rest/internal/domain"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	userHandler *userHandler,
	noteHandler *noteHandler,
	middlewares *middlewares,
) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, domain.Errorf(domain.NOTFOUND_ERROR, "resource not found"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		writeHTTPError(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	r.Use(
		middlewares.RequestID,
		middlewares.Logger,
		middlewares.Recovery,
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/user/register", userHandler.register)
			r.Post("/user/login", userHandler.login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth)

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
