package main

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/register/", app.home)
	mux.HandleFunc("/login/", app.home)
	mux.HandleFunc("/logout/", app.home)
	mux.HandleFunc("/forget-password/", app.home)
	mux.HandleFunc("/user-info/", app.home)
	mux.HandleFunc("/update/user-info/", app.home)

	mux.HandleFunc("/add/product/id/", app.home)
	mux.HandleFunc("/delete/product/id/", app.home)
	mux.HandleFunc("/update/product/id/", app.home)
	mux.HandleFunc("/price/product/id/", app.home)

	// mux.HandleFunc("/order/cart/id/", app.home)
	mux.HandleFunc("/order/user-info/id/", app.home)
	mux.HandleFunc("/order/payment/id/", app.home)
	mux.HandleFunc("/order/publish/id/", app.home)

	mux.HandleFunc("/listing/", app.home)
	mux.HandleFunc("/listing/advance-pattern/", app.home)

	mux.HandleFunc("/product/info/id/", app.home)

	return mux

}
