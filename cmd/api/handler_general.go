package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	fmt.Fprint(w, "Site is working")
}

/* to print message */
func (app *Application) Message(w http.ResponseWriter, statusCode int, key, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{key: value})
}

func (app *Application) ErrorMessage(w http.ResponseWriter, statusCode int, Message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Status bool
		Error  string
	}{
		Status: false,
		Error:  Message,
	})
}

func (app *Application) sendJSONResponse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
