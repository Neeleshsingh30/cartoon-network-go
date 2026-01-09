package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type OmdbResponse struct {
	IMDBRating string `json:"imdbRating"`
}

func FetchIMDBRating(title string) (float32, error) {
	apiKey := os.Getenv("OMDB_API_KEY")

	endpoint := fmt.Sprintf(
		"http://www.omdbapi.com/?t=%s&apikey=%s",
		url.QueryEscape(title),
		apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data OmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	rating, err := strconv.ParseFloat(data.IMDBRating, 32)
	if err != nil {
		return 0, err
	}

	return float32(rating), nil
}
