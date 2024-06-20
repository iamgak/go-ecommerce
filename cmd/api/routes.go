package main

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/api/register/", app.UserRegister)
	mux.HandleFunc("/api/login/", app.UserLogin)
	mux.HandleFunc("/api/activation_token/", app.UserActivationToken)
	mux.HandleFunc("/api/logout/", app.UserLogout)
	mux.HandleFunc("/api/forget-password/", app.UserForgetPassword)
	mux.HandleFunc("/api/user-info/", app.home)
	mux.HandleFunc("/api/update/user-info/", app.home)

	mux.HandleFunc("/api/add/product/", app.CreateObject)
	mux.HandleFunc("/api/delete/product/id/", app.DeleteObject)
	mux.HandleFunc("/api/update/product/id/", app.CreateObject)
	mux.HandleFunc("/api/price/product/id/", app.home)
	mux.HandleFunc("/api/add/cart/", app.CreateCart)
	mux.HandleFunc("/api/order/user-info/id/", app.home)
	mux.HandleFunc("/api/order/payment/id/", app.home)
	mux.HandleFunc("/api/order/publish/id/", app.OrderReview)
	mux.HandleFunc("/api/order/cancel/id/", app.CancelOrder)

	mux.HandleFunc("/api/listing/", app.home)
	mux.HandleFunc("/api/listing/advance-pattern/", app.home)

	mux.HandleFunc("/api/product/info/id/", app.home)

	return mux

}
