package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Post("/get", GetEntry)
	r.Post("/update", UpdateEntry)
	r.Post("/create", CreateEntry)
	log.Info().Msg("Setting up routes")
	return r
}
