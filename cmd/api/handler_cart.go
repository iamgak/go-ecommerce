package main

import (
	"encoding/json"
	data "github.com/iamgak/go-ecommerce/internals"
	"net/http"
	"strconv"
)

func (app *Application) CreateCart(w http.ResponseWriter, r *http.Request) {
	var input *data.Cart
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON payload", http.StatusInternalServerError)
		return
	}

	validator := app.Model.Carts.ErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Model.Carts.AddInCart(input)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Added in cart")
}

func (app *Application) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	cart_id, err := strconv.Atoi(r.URL.Query().Get("cart_id"))
	if err != nil {
		app.NotFound(w)
		app.ErrorLog.Print(err)
		return
	}

	err = app.Model.Carts.RemoveFromCart(cart_id, app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Removed From Cart")
}

func (app *Application) CartListing(w http.ResponseWriter, r *http.Request) {
	data, err := app.Model.Carts.CartListing(app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, data)
}
