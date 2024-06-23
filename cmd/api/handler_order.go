package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) CancelOrder(w http.ResponseWriter, r *http.Request) {
	uid := 1
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.Order.CancelOrder(uid, order_id)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Order Cancelled")
}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input *models.Order
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

	validator := app.Order.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.Order.CreateOrder(input)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) MakePayment(w http.ResponseWriter, r *http.Request) {
	var input *models.Payment
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

	validator := app.Payment.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	// err = app.Order.CreateOrder(input)
	// if err != nil {
	// 	app.Message(w, 500, "Internal Error", "Internal Server Error")
	// 	return
	// }

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) OrderReview(w http.ResponseWriter, r *http.Request) {
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if order_id == 0 || err != nil {
		http.NotFound(w, r)
		return
	}

}
