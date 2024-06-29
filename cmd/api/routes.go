package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	//home page
	mux.HandleFunc("/", app.home)
	// customer/user
	mux.HandleFunc("/api/user/register/", app.UserRegister)
	mux.HandleFunc("/api/user/login/", app.UserLogin)
	mux.HandleFunc("/api/user/activate_account/", app.UserActivationToken)
	mux.HandleFunc("/api/user/logout/", app.UserLogout)
	mux.HandleFunc("/api/user/forget-password/", app.UserForgetPassword)
	mux.HandleFunc("/api/user/new-password/", app.NewPassword)

	//seller
	mux.HandleFunc("/api/seller/register/", app.SellerRegister)
	mux.HandleFunc("/api/seller/login/", app.SellerLogin)
	mux.HandleFunc("/api/seller/forget-password/", app.SellerForgetPassword)
	mux.HandleFunc("/api/seller/new-password/", app.SellerNewPassword)
	mux.HandleFunc("/api/seller/logout/", app.SellerLogout)
	mux.HandleFunc("/api/seller/activate_account/", app.SellerActivationToken)

	seller := alice.New(app.SellerAuthentication)

	mux.Handle("/api/add/product/", seller.ThenFunc(app.CreateProduct))
	mux.Handle("/api/add/product/addr/", seller.ThenFunc(app.CreateProductAddr))
	mux.Handle("/api/delete/product/", seller.ThenFunc(app.DeleteProduct))
	mux.Handle("/api/update/product/", seller.ThenFunc(app.UpdateProduct))

	customer := alice.New(app.UserAuthentication)

	mux.Handle("/api/add/cart/", customer.ThenFunc(app.CreateCart))
	mux.Handle("/api/add/order/", customer.ThenFunc(app.CreateOrder))
	mux.Handle("/api/order/payment/", customer.ThenFunc(app.MakePayment))
	mux.Handle("/api/order/review/", customer.ThenFunc(app.OrderReview))
	mux.Handle("/api/order/cancel/", customer.ThenFunc(app.CancelOrder))

	mux.HandleFunc("/api/listing/", app.home)
	mux.HandleFunc("/api/product/view/", app.home)
	// mux.HandleFunc("/api/listing/", app.home)

	standard := alice.New(app.logRequest, secureHeaders)
	return standard.Then(mux)

}
