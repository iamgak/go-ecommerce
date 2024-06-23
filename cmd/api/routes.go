package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/api/register/", app.UserRegister)
	mux.HandleFunc("/api/seller/register/", app.SellerRegister)
	mux.HandleFunc("/api/user/login/", app.UserLogin)
	mux.HandleFunc("/api/user/activation_token/", app.UserActivationToken)
	mux.HandleFunc("/api/user/logout/", app.UserLogout)
	mux.HandleFunc("/api/user/forget-password/", app.UserForgetPassword)

	loggedIn := alice.New(app.authenticate, app.MyMiddlerware)
	mux.Handle("/api/add/product/", loggedIn.ThenFunc(app.CreateObject))
	mux.Handle("/api/delete/product/", loggedIn.ThenFunc(app.DeleteObject))
	mux.Handle("/api/update/product/", loggedIn.ThenFunc(app.UpdateObject))
	mux.Handle("/api/add/cart/", loggedIn.ThenFunc(app.CreateCart))
	mux.Handle("/api/add/order/", loggedIn.ThenFunc(app.CreateOrder))

	mux.Handle("/api/order/payment/", loggedIn.ThenFunc(app.MakePayment))
	mux.Handle("/api/order/review/", loggedIn.ThenFunc(app.OrderReview))
	mux.Handle("/api/order/cancel/", loggedIn.ThenFunc(app.CancelOrder))

	mux.HandleFunc("/api/listing/", app.home)
	// mux.HandleFunc("/api/listing/", app.home)

	mux.HandleFunc("/api/product/view/", app.home)

	standard := alice.New(app.logRequest, secureHeaders)
	return standard.Then(mux)

}
