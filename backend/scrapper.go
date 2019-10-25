package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aosousa/go-movielookup/models"
	utils "github.com/aosousa/golang-utils"
)

func getExcelDocumentRows(worksheet string) [][]string {
	XLSXLocation := "Movies.xlsx"

	// open XLSX document
	xlsx, err := excelize.OpenFile(XLSXLocation)
	if err != nil {
		utils.HandleError(err)
	}

	// get all rows in worksheet specified
	rows := xlsx.GetRows(worksheet)

	return rows
}

func buildMoviesWatchlistSlice(rows [][]string) [27][]string {
	var movies [27][]string

	for _, row := range rows {
		for i := 0; i < 27; i++ {
			if row[i] != "" {
				movies[i] = append(movies[i], row[i])
			}
		}
	}

	return movies
}

func buildMoviesWatchedSlice(rows [][]string) []string {
	var movies []string

	for _, row := range rows {
		if row[0] != "" {
			movies = append(movies, row[0])
		}
	}

	return movies
}

func getMovieOMDBData(queryURL, title string) []byte {
	res, err := http.Get(queryURL)
	if res.StatusCode != 200 || err != nil {
		utils.HandleError(err)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.HandleError(err)
	}

	return content
}

func saveWatchlistAsJSON(column int) {
	var apiError models.Error

	rows := getExcelDocumentRows("Movies")
	movieSlice := buildMoviesWatchlistSlice(rows)
	file, movies := openJSONFile("watchlist.json")
	defer file.Close()

	for row := range movieSlice[column] {
		// build OMDB query
		title := strings.Replace(movieSlice[column][row], " ", "+", -1)
		queryURL := fmt.Sprintf("%st=%s&type=movie", baseURL, title)

		// call OMDB API
		movieBytes := getMovieOMDBData(queryURL, title)
		json.Unmarshal(movieBytes, &apiError)
		if apiError.Response == "True" {
			fmt.Printf("Adding movie: %s\n", movieSlice[column][row])

			var movie models.Movie
			json.Unmarshal(movieBytes, &movie)

			movies = append(movies, movie)
		} else {
			fmt.Printf("Error occurred during OMDB request for movie: %s\n", movieSlice[column][row])
		}
	}

	updateJSONFile(file, movies)
}

func saveMoviesWatchedAsJSON(start int) {
	var apiError models.Error

	rows := getExcelDocumentRows("movie ratings")
	movieSlice := buildMoviesWatchedSlice(rows)
	file, movies := openJSONFile("watched.json")
	defer file.Close()

	for i := start; i < start+50; i++ {
		// build OMDB query
		title := strings.Replace(movieSlice[i], " ", "+", -1)
		queryURL := fmt.Sprintf("%st=%s&type=movie", baseURL, title)

		// call OMDB API
		movieBytes := getMovieOMDBData(queryURL, title)
		json.Unmarshal(movieBytes, &apiError)
		if apiError.Response == "True" {
			fmt.Printf("Adding movie: %s\n", movieSlice[i])

			var movie models.Movie
			json.Unmarshal(movieBytes, &movie)

			movies = append(movies, movie)
		} else {
			fmt.Printf("Error occurred during OMDB request for movie: %s\n", movieSlice[i])
		}
	}

	updateJSONFile(file, movies)
}

func openJSONFile(filename string) (*os.File, []models.Movie) {
	var movies []models.Movie

	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		utils.HandleError(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		utils.HandleError(err)
	}

	err = json.Unmarshal(bytes, &movies)
	if err != nil {
		utils.HandleError(err)
	}

	return file, movies
}

func updateJSONFile(file *os.File, movies []models.Movie) {
	newBytes, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
		utils.HandleError(err)
	}

	_, err = file.WriteAt(newBytes, 0)
	if err != nil {
		utils.HandleError(err)
	}
}
