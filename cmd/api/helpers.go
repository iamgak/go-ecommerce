package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	models "github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) MyMiddlerware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.ValidUser(r)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		if userID > 0 {
			app.isAuthenticated = true
		}

		if !app.isAuthenticated {
			http.Redirect(w, r, "/api/user/login/", http.StatusSeeOther)
			return
		}

		app.Uid = userID
		next.ServeHTTP(w, r)
	})
}

func (app *Application) UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.ValidUser(r)
		if err != nil && err != models.ErrNoCookieFound {
			app.ServerError(w, err)
			return
		}

		if userID <= 0 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing authentication token"
			app.ErrorMessage(w, http.StatusUnauthorized, message)
			return
		}

		app.isAuthenticated = true
		app.Uid = userID
		app.InfoLog.Printf("userId %d online", app.Uid)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) SellerAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.ValidSeller(r)
		if err != nil && err != models.ErrNoCookieFound {
			app.ServerError(w, err)
			return
		}

		if userID <= 0 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			message := "invalid or missing authentication token"
			app.ErrorMessage(w, http.StatusUnauthorized, message)
			return
		}

		// Authentication successful
		app.isAuthenticated = true
		app.Uid = userID
		app.InfoLog.Printf("sellerId %d online", app.Uid)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	app.ErrorMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) CustomError(w http.ResponseWriter, err error, message string, status int) {
	app.ErrorMessage(w, status, message)
	app.ErrorLog.Print(err)
}
