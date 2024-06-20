package main

import (
	"database/sql"
	"net/http"
)

func (app *Application) NoRecordFound(err error) (bool, error) {
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

func (app *Application) MyMiddlerware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id, authenticated := app.ValidUser(r)
		if !authenticated {
			http.Redirect(w, r, "/api/login/", http.StatusSeeOther)
			return
		}
		app.Authenticated = authenticated
		app.Uid = user_id
		// app.seller = false
		next.ServeHTTP(w, r)
	})
}
