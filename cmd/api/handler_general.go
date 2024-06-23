package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Site is working")
}

func (app *Application) ValidUser(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("ldata")
	if err != nil || cookie.Value == "" {
		return 0, false
	}

	id, seller := app.User.ValidToken(cookie.Value)
	return id, seller
}

/* to print message */
func (app *Application) Message(w http.ResponseWriter, statusCode int, key, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{key: value})
}

func (app *Application) sendJSONResponse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
