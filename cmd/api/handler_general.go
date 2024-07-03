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

func (app *Application) ErrorMessage(w http.ResponseWriter, statusCode int, Message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Status bool
		Error  any
	}{
		Status: false,
		Error:  Message,
	})
}

func (app *Application) FinalMessage(w http.ResponseWriter, statusCode int, Message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Success bool
		Message any
	}{
		Success: true,
		Message: Message,
	})
}
