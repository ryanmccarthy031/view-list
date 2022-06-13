package routes

import (
	"github.com/go-chi/chi/v5"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Post("/get", GetEntry)
	r.Post("/update", UpdateEntry)
	r.Post("/create", CreateEntry)

	return r
}
