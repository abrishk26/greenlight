package main

import (
	"fmt"
	"net/http"
)

func (a *application) logError(_ *http.Request, err error) {
	a.logger.Println(err)
}

func (a *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := map[string]interface{}{
		"error": message,
	}

	err := a.writeJSON(w, env, nil, status)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	a.errorResponse(w, r, 500, message)
}

func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the resource could not be found"
	a.errorResponse(w, r, http.StatusNotFound, message)
}

func (a *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this request", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
