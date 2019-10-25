package main

import (
	"os"
	"strconv"

	"github.com/aosousa/go-movielookup/models"
)

var (
	baseURL string
	config  models.Config
)

func init() {
	config = models.CreateConfig()
	baseURL = "http://www.omdbapi.com/?apiKey=" + config.APIKey + "&"
}

func main() {
	args := os.Args
	cmd := args[1]
	value, _ := strconv.Atoi(os.Args[2])

	switch cmd {
	case "--watchlist":
		saveWatchlistAsJSON(value)
	case "--watched":
		saveMoviesWatchedAsJSON(value)
	}
}
