package main

import (
	"net/http"
)

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": a.conf.env,
		"version":     version,
	}

	err := a.writeJSON(w, data, nil, 200)

	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}
