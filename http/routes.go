package http

import (
	"companyXYZ-waitlistAPI/env"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	companyXYZHOMEURL string = "https://www.companyXYZ.xyz/"
	Production        string = "production"
)

func (a Application) Routes() *chi.Mux {

	r := chi.NewRouter()

	if env.GetEnvVar().Server.Env != Production {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)

	if env.GetEnvVar().Server.Env == Production {
		r.Use(httprate.LimitByIP(5, 10*time.Minute))

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
	r.Get("/healthcheck", a.healthcheckHandler)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:"+env.GetEnvVar().Server.Port+"/swagger/doc.json"),
	))

	return r
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	return origin == companyXYZHOMEURL
}
