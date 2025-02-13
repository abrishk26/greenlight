package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating movie...")
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	
	fmt.Fprintf(w, "Showing details of movies with id: %d", id)
}

