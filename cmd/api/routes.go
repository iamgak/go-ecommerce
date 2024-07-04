package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	// router := http.NewServeMux()
	router := httprouter.New()

	//404 page
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	//home page
	router.HandlerFunc(http.MethodGet, "/", app.home)

	// customer
	router.HandlerFunc(http.MethodPost, "/api/user/register/", app.UserRegister)
	router.HandlerFunc(http.MethodPost, "/api/user/login/", app.UserLogin)
	router.HandlerFunc(http.MethodGet, "/api/user/activate_account/", app.UserActivationToken)
	router.HandlerFunc(http.MethodGet, "/api/user/logout/", app.UserLogout)
	router.HandlerFunc(http.MethodPost, "/api/user/forget-password/", app.UserForgetPassword)
	router.HandlerFunc(http.MethodPost, "/api/user/new-password/", app.NewPassword)

	//seller
	router.HandlerFunc(http.MethodPost, "/api/seller/register/", app.SellerRegister)
	router.HandlerFunc(http.MethodPost, "/api/seller/login/", app.SellerLogin)
	router.HandlerFunc(http.MethodPost, "/api/seller/forget-password/", app.SellerForgetPassword)
	router.HandlerFunc(http.MethodPost, "/api/seller/new-password/", app.SellerNewPassword)
	router.HandlerFunc(http.MethodGet, "/api/seller/logout/", app.SellerLogout)
	router.HandlerFunc(http.MethodGet, "/api/seller/activate_account/", app.SellerActivationToken)

	seller := alice.New(app.SellerAuthentication)

	router.Handler(http.MethodPost, "/api/add/product/", seller.ThenFunc(app.CreateProduct))
	router.Handler(http.MethodPost, "/api/add/product/addr/", seller.ThenFunc(app.CreateProductAddr))
	router.Handler(http.MethodPut, "/api/delete/product/", seller.ThenFunc(app.DeleteProduct))
	router.Handler(http.MethodPut, "/api/update/product/", seller.ThenFunc(app.UpdateProduct))

	customer := alice.New(app.UserAuthentication)

	router.Handler(http.MethodPost, "/api/add/cart/", customer.ThenFunc(app.CreateCart))
	router.Handler(http.MethodDelete, "/api/cart/delete/", customer.ThenFunc(app.RemoveFromCart))
	router.Handler(http.MethodGet, "/api/cart/listing/", customer.ThenFunc(app.CartListing))
	router.Handler(http.MethodPost, "/api/add/order/", customer.ThenFunc(app.CreateOrder))

	router.Handler(http.MethodGet, "/api/order/listing/", customer.ThenFunc(app.OrderListing))
	router.Handler(http.MethodPost, "/api/order/payment/", customer.ThenFunc(app.MakePayment))
	router.Handler(http.MethodGet, "/api/order/review/", customer.ThenFunc(app.OrderReview))
	router.Handler(http.MethodPatch, "/api/order/cancel/", customer.ThenFunc(app.CancelOrder))

	router.HandlerFunc(http.MethodGet, "/api/product/listing/", app.ProductListing)
	router.Handler(http.MethodGet, "/api/dbg/", expvar.Handler())

	//Added RateLimiter, SecureHeader, RecoverPanic, running status / metrics
	standard := alice.New(app.rateLimit, secureHeaders, app.recoverPanic, app.metrics)

	return standard.Then(router)

}
