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

	var input struct {
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.editConflictResponse(w, r)
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, movie, nil, http.StatusOK)
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

	err = a.writeJSON(w, map[string]string{"message": "movie successfully deleted"}, nil, http.StatusOK)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Genres []string
		Filter data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = a.readString(qs, "title", "")
	input.Genres = a.readCSV(qs, "genres", []string{})
	input.Filter.Page = a.readInt(qs, "page", 1, v)
	input.Filter.PageSize = a.readInt(qs, "page_size", 20, v)
	input.Filter.Sort = a.readString(qs, "sort", "id")

	input.Filter.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-year", "-runtime", "-title"}

	if data.ValidateFilters(v, input.Filter); !v.Valid(){
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := a.models.Movies.GetAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, map[string]interface{}{ "movies": movies, "metadata": metadata }, nil, http.StatusOK)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}
