package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/rs/zerolog/log"
)

type getReq struct {
	Title string `json:"title,omitempty"`
	Year  string `json:"release_date,omitempty"`
	Type  string `json:"media_type,omitempty"`
}

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal().Err(err).Msg("failed to process environment variables")
	}
	return os.Getenv(key)
}

func GetEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Get Entry")
	client := notionapi.NewClient(notionapi.Token(GoDotEnvVariable("NOTION_KEY")))

	var item getReq

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	year, err := strconv.ParseFloat(item.Year, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse year")
	}

	types := map[string]string{"tv": "TV Series", "movie": "Movie"}

	req := &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			&notionapi.PropertyFilter{
				Property: "Title",
				RichText: &notionapi.TextFilterCondition{
					Equals: item.Title,
				},
			},
			notionapi.PropertyFilter{
				Property: "Release Year",
				Number: &notionapi.NumberFilterCondition{
					Equals: &year,
				},
			},
			notionapi.PropertyFilter{
				Property: "Type",
				Select: &notionapi.SelectFilterCondition{
					Equals: types[item.Type],
				},
			},
		},
	}

	results, err := client.Database.Query(context.Background(), notionapi.DatabaseID(GoDotEnvVariable("NOTION_DATABASE")), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve Database Results")
	}
	json.NewEncoder(w).Encode(results.Results)

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
