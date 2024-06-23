package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) CreateCart(w http.ResponseWriter, r *http.Request) {
	var input *models.Cart
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

	app.Message(w, 200, "Message", "Added in Cart")
}
