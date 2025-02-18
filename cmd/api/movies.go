package main

import (
	"errors"
	"fmt"
	"net/http"

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
		Title:   input.Title,
		Runtime: input.Runtime,
		Year:    input.Year,
		Genres:  input.Genres,
	}

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Movies.Insert(movie)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = a.writeJSON(w, movie, headers, http.StatusCreated)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	// fmt.Fprintf(w, "%+v\n", input)
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)

	if err != nil || id < 1 {
		a.errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	movie, err := a.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, movie, nil, 200)

	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

}

func (a *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	movie := &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	movie.ID = id
	err = a.models.Movies.Update(movie)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = a.writeJSON(w, movie, headers, http.StatusAccepted)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	err = a.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, map[string]string{"message": "movie successfully deleted"}, nil , http.StatusOK)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
