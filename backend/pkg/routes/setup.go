package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Post("/api/get", GetEntry)
	r.Get("/api/ping", WakeUp)
	r.Post("/api/update", UpdateEntry)
	r.Post("/api/create", CreateEntry)
	log.Info().Msg("Setting up routes")
	return r
}
