package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", a.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", a.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", a.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)

	return a.recoverPanic(a.rateLimit(router))
}
