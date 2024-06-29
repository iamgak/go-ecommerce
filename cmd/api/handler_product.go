package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.NotFound(w, r)
		app.InfoLog.Print(err)
		return
	}

	var input *models.Product
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.Product.ProductErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Product.UpdateProduct(input, product_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "Product Updated")
}

func (app *Application) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input *models.Product
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.Product.ProductErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Product.CreateProduct(input)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "Product Added")
}

func (app *Application) CreateProductAddr(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if product_id == 0 || err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	err = app.Product.UserProductExist(product_id, app.Uid)
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	var input *models.Product_Addr
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		message := ErrResp{Error: "Incorrect Format Data"}
		app.sendJSONResponse(w, 200, message)
		return
	}

	validator := app.Product.ProductAddrErrorCheck(input)
	if len(validator) != 0 {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	input.Order_id = product_id
	err = app.Product.CreateProductAddr(input)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "Address Added to Product")
}

func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.NotFound(w, r)
		app.InfoLog.Print(err)
		return
	}

	err = app.Product.DeleteProduct(app.Uid, product_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "Message", "Product Deleted")
}

func (app *Application) UpdateProductQuantity(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	quantity, err1 := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil || err1 != nil {
		app.ServerError(w, err)
		return
	}

	if product_id == 0 {
		app.NotFound(w)
		return
	}

	err = app.Product.UpdateProductQuantity(product_id, app.Uid, quantity)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.Message(w, 200, "message", "Product Quantity updated")
}
