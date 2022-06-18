package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/rs/zerolog/log"
)

type getReq struct {
	Title string `json:"title,omitempty"`
	Year  int    `json:"release_date,omitempty"`
	Type  string `json:"media_type,omitempty"`
}
type updateReq struct {
	ID   notionapi.PageID            `json:"id,omitempty"`
	Page notionapi.PageUpdateRequest `json:"page"`
}
type createReq struct {
	Page notionapi.PageCreateRequest `json:"page"`
}

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal().Err(err).Msg("failed to process environment variables")
	}
	return os.Getenv(key)
}

func WakeUp(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Waking Server")
}

func GetEntry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Called Get Entry")
	client := notionapi.NewClient(notionapi.Token(GoDotEnvVariable("NOTION_KEY")))

	var item getReq

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	year := float64(item.Year)
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
	client := notionapi.NewClient(notionapi.Token(GoDotEnvVariable("NOTION_KEY")))
	log.Info().Msg("Called Update Entry")
	var item updateReq
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	results, err := client.Page.Update(context.Background(), notionapi.PageID(item.ID), &item.Page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update Page")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results.ID)
}
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	client := notionapi.NewClient(notionapi.Token(GoDotEnvVariable("NOTION_KEY")))
	log.Info().Msg("Called Update Entry")
	var item createReq
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	item.Page.Parent = notionapi.Parent{
		Type:       "database_id",
		DatabaseID: notionapi.DatabaseID(GoDotEnvVariable("NOTION_DATABASE")),
	}
	results, err := client.Page.Create(context.Background(), &item.Page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update Page")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}
