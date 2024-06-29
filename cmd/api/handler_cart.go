package main

import (
	"encoding/json"
	models "github.com/iamgak/go-ecommerce/internals"
	"net/http"
)

func (app *Application) CreateCart(w http.ResponseWriter, r *http.Request) {
	var input *models.Cart
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON payload", http.StatusInternalServerError)
		return
	}

	validator := app.Cart.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Cart.AddInCart(input)
	if err != nil {
		app.CustomError(w, err, "Internal Server error", http.StatusInternalServerError)
		return
	}

	app.Message(w, 200, "Message", "Added in Cart")
}
