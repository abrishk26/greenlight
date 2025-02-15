package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/abrishk26/greenlight/internal/data"
	"github.com/abrishk26/greenlight/internal/validator"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := a.readJSON(w, r, &input)

	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	movie := &data.Movie{
		Title: input.Title,
		Runtime: input.Runtime,
		Year: input.Year,
		Genres: input.Genres,
	}
	
	data.ValidateMovie(v, movie)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)

	if err != nil || id < 1 {
		a.errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   132,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = a.writeJSON(w, movie, nil, 200)

	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}
