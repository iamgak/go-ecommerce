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
		app.Uid, app.Seller = app.ValidUser(r)
		if app.Uid == 0 {
			http.Redirect(w, r, "/api/user/login/", http.StatusSeeOther)
			// http.NotFound(w,r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func Adapt(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
