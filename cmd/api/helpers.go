package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type ErrResp struct {
	Error any
}

type SuccessResp struct {
	Status  bool
	Success any
}

func (app *Application) MyMiddlerware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id, err := app.ValidUser(r)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		if user_id > 0 {
			app.isAuthenticated = true
		}

		if !app.isAuthenticated {
			http.Redirect(w, r, "/api/user/login/", http.StatusSeeOther)
			return
		}

		app.Uid = user_id
		next.ServeHTTP(w, r)
	})
}

func (app *Application) UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id, err := app.ValidUser(r)
		if err != nil {
			app.NotFound(w)
			app.ErrorLog.Print(err)
			return
		}

		if user_id <= 0 {
			http.Redirect(w, r, "/api/user/login/", http.StatusSeeOther)
			return
		}

		app.isAuthenticated = true
		app.Uid = user_id
		app.InfoLog.Printf("userId %d online", app.Uid)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) SellerAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.ValidSeller(r)
		if err != nil {
			app.InfoLog.Print(err)
			app.NotFound(w)
			return
		}

		if userID <= 0 {
			http.Redirect(w, r, "/api/seller/login/", http.StatusSeeOther)
			return
		}

		// Authentication successful
		app.isAuthenticated = true
		app.Uid = userID
		app.InfoLog.Printf("sellerId %d online", app.Uid)
		next.ServeHTTP(w, r)
	})
}

func Adapt(handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	app.ErrorMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) CustomError(w http.ResponseWriter, err error, message string, status int) {
	app.ErrorMessage(w, status, message)
	app.ErrorLog.Print(err)
}
