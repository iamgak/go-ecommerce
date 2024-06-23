package main

import (
	"encoding/json"
	"github.com/iamgak/go-ecommerce/internals"
	"net/http"
	"strconv"
)

func (app *Application) UpdateObject(w http.ResponseWriter, r *http.Request) {
	object_id, err := strconv.Atoi(r.URL.Query().Get("object_id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var input *models.Object
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.Message(w, 200, "Invalid", "Data Incorrect format")
		return
	}

	validator := app.Object.ErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Object.UpdateObject(input, object_id)
	if err != nil {
		app.Message(w, 500, "Internal Error", "Internal Server Error")
		return
	}

	app.Message(w, 200, "Message", "Product Updated")
}

func (app *Application) CreateObject(w http.ResponseWriter, r *http.Request) {
	var input *models.Object
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

	input.Uid = app.Uid
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
