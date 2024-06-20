package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/iamgak/go-ecommerce/internals/cart"
	"github.com/iamgak/go-ecommerce/internals/checkout"
	"github.com/iamgak/go-ecommerce/internals/object"
	"github.com/iamgak/go-ecommerce/internals/user"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Site is working")
}

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request) {
	var input *user.User
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	token, err := app.User.UserLogin(input)
	if err != nil {
		if err == sql.ErrNoRows {
			app.Message(w, 200, "Invalid", "Incorrect Credentails")
			return
		}

		app.Message(w, 200, "Internal Error", "Internal Server Error")
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	cookie := &http.Cookie{
		Name:    "ldata",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	app.Message(w, 200, "Message", "Login Successfull")
}

func (app *Application) UserRegister(w http.ResponseWriter, r *http.Request) {
	var input *user.User
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.User.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	inValid, err := app.User.EmailExist(input.Email)
	if err != nil && err != sql.ErrNoRows {
		app.InfoLog.Print(err)
		return
	}

	if inValid {
		app.Message(w, 200, "Email", "Email already Exist")
		return
	}

	err = app.User.CreateUser(input)
	if err != nil {
		app.InfoLog.Print(err)
		return
	}

	app.sendJSONResponse(w, 200, input)
}

func (app *Application) UserActivationToken(w http.ResponseWriter, r *http.Request) {
	activation_token := r.URL.Query().Get("activation_token")
	if r.Method != "GET" || activation_token == "" {
		http.NotFound(w, r)
		return
	}

	err := app.User.ActivateUser(activation_token)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Your Account has been Verified")
}

func (app *Application) UserForgetPassword(w http.ResponseWriter, r *http.Request) {
}

func (app *Application) UserLogout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	// app.InfoLog.Panic(authHeader)
	if authHeader == "" {
		http.NotFound(w, r)
		return
	}

	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		http.NotFound(w, r)
		return
	}

	err := app.User.UserLogout(authHeader)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 500, "Message", "Logout Successfully")
}

func (app *Application) CreateObject(w http.ResponseWriter, r *http.Request) {
	var input *object.Object
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.Object.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = 1
	err = app.Object.CreateObject(input)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) DeleteObject(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.Object.DeleteObject(app.Uid, product_id)
	if err != nil {
		app.Message(w, 500, "Server Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Product Deleted")
}

func (app *Application) CreateCart(w http.ResponseWriter, r *http.Request) {
	var input *cart.Cart
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	input.ProductId = product_id
	validator := app.Cart.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = 1
	err = app.Cart.AddInCart(input)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) CancelOrder(w http.ResponseWriter, r *http.Request) {
	uid := 1
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.Checkout.CancelOrder(uid, order_id)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Order Cancelled")
}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input *checkout.Order
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	cartId, err := strconv.Atoi(r.URL.Query().Get("cartId"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	input.CartId = cartId

	validator := app.Checkout.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	// err = app.Cart.AddInCart(input)
	// if err != nil {
	// 	app.Message(w, 500, "Internal Error", "Internal Server Error")
	// 	return
	// }

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) OrderReview(w http.ResponseWriter, r *http.Request) {
}

func (app *Application) Message(w http.ResponseWriter, statusCode int, key, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{key: value})
}

// to print json message
func (app *Application) sendJSONResponse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

func (app *Application) ValidUser(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("ldata")
	if err != nil || cookie.Value == "" {
		return 0, false
	}

	id, _ := app.User.ValidToken(cookie.Value)
	return id, true
}
