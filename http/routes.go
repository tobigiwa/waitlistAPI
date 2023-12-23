package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a Application) Routes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/waitlist", a.waitListHandler)
	r.Post("/confirm", a.confirmAndSaveHandler)

	return r

}
