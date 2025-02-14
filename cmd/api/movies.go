package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abrishk26/greenlight/internal/data"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating movie...")
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)

	if err != nil || id < 1 {
		a.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	
	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Runtime: 132,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}
	
	err = a.writeJSON(w, movie, nil, 200)
	
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}
