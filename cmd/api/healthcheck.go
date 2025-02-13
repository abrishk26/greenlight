package main

import (
	"fmt"
	"net/http"
)

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "status: Available")
	fmt.Fprintf(w, "environment: %s\n", a.conf.env)
	fmt.Fprintf(w, "listening on port: %d\n", a.conf.port)
}