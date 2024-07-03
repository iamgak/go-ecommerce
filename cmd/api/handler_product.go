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
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	var input *models.Product
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON Payload", http.StatusInternalServerError)
		return
	}

	validator := app.Model.Products.ProductErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Model.Products.UpdateProduct(input, product_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Product Updated")
}

func (app *Application) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input *models.Product
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.CustomError(w, err, "Invalid Json Payload", http.StatusInternalServerError)
		return
	}

	validator := app.Model.Products.ProductErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	input.Uid = app.Uid
	err = app.Model.Products.CreateProduct(input)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Product Added")
}

func (app *Application) CreateProductAddr(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if product_id == 0 || err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	err = app.Model.Products.UserProductExist(product_id, app.Uid)
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	var input *models.Product_Addr
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON Payload", http.StatusInternalServerError)
		return
	}

	validator := app.Model.Products.ProductAddrErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	input.Order_id = product_id
	err = app.Model.Products.CreateProductAddr(input)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Seller Address Added")
}

func (app *Application) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id, err := strconv.Atoi(r.URL.Query().Get("product_id"))
	if err != nil {
		http.NotFound(w, r)
		app.InfoLog.Print(err)
		return
	}

	err = app.Model.Products.DeleteProduct(app.Uid, product_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Product Deleted")
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

	err = app.Model.Products.UpdateProductQuantity(product_id, app.Uid, quantity)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Product Quantity Updated")
}

func (app *Application) ProductListing(w http.ResponseWriter, r *http.Request) {
	data, err := app.Model.Products.ProductListing()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, data)
}
