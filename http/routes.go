package http

import (
	"Blockride-waitlistAPI/env"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

const (
	blockrideHOMEURL string = "https://www.blockride.xyz/"
	Production       string = "production"
)

func (a Application) Routes() *chi.Mux {

	r := chi.NewRouter()

	if env.GetEnvVar().Server.Env != Production {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)

	if env.GetEnvVar().Server.Env == Production {
		r.Use(httprate.LimitByIP(3, 10*time.Minute))

		r.Use(cors.Handler(
			cors.Options{
				AllowOriginFunc:  AllowOriginFunc,
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Accept", "Content-Type"},
				AllowCredentials: true,
			}))
	}

	r.Post("/joinwaitlist", a.waitListHandler)
	r.Get("/confirmuser", a.confirmAndSaveHandler)

	return r
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	return origin == blockrideHOMEURL
}
