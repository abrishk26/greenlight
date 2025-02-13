package main

import (
	"fmt"
	"net/http"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating movie...")
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Showing details of movies with id: %d", id)
}
