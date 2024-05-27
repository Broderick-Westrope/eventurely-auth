package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func setupRouter(websiteDomain string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{websiteDomain},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: append([]string{"Content-Type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	// SuperTokens Middleware
	r.Use(supertokens.Middleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	return r
}
