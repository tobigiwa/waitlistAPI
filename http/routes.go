package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a Application) Routes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/joinwaitlist", a.waitListHandler)
	r.Get("/confirmuser", a.confirmAndSaveHandler)

	return r
}
