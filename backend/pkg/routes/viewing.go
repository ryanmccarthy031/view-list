package routes

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Get Entry")
}
func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Update Entry")
}
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Create Entry")
}
