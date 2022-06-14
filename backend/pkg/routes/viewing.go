package routes

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Get Entry")
	t := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Getting Entry")),
	}

	buff := bytes.NewBuffer(nil)
	t.Write(buff)
}
func UpdateEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Update Entry")
	t := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Updating Entry")),
	}

	buff := bytes.NewBuffer(nil)
	t.Write(buff)
}
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Create Entry")
	t := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("Creating Entry")),
	}

	buff := bytes.NewBuffer(nil)
	t.Write(buff)
}
