package main

import "github.com/aosousa/go-movielookup/models"

var (
	baseURL string
	config  models.Config
)

func init() {
	config = models.CreateConfig()
	baseURL = "http://www.omdbapi.com/?apiKey=" + config.APIKey + "&"
}

func main() {
	// saveWatchlistAsJSON(2)
	saveMoviesWatchedAsJSON()
}
