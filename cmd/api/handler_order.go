package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	models "github.com/iamgak/go-ecommerce/internals"
)

func (app *Application) OrderListing(w http.ResponseWriter, r *http.Request) {
	order_listing, err := app.Order.OrderListing(app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, order_listing)
}

func (app *Application) CancelOrder(w http.ResponseWriter, r *http.Request) {
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	err, active := app.Order.OrderStatus(order_id, app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if active {
		err = app.Order.CancelOrder(app.Uid, order_id)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}

	// app.Message(w, 200, "Message", "Order Cancelled")
	http.Redirect(w, r, fmt.Sprintf("/api/order/review/?order_id=%d", order_id), http.StatusPermanentRedirect)
}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.CustomError(w, err, "Fail to load request body", http.StatusInternalServerError)
		return
	}

	// Decode the first struct
	requestData := &models.RequestData{}
	err = json.Unmarshal(body, requestData)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON payload", http.StatusInternalServerError)
		return
	}

	// Decode the second struct
	input := &models.UserAddr{}
	err = json.Unmarshal(body, input)
	if err != nil {
		app.CustomError(w, err, "Invalid JSON payload", http.StatusInternalServerError)
		return
	}

	requestData.UserId = app.Uid
	validator := app.Order.RequestErrorCheck(requestData)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	product_id, product_quantity, required_quantity, price, err := app.Order.ValidCart(requestData.CartID, app.Uid)
	if err != nil {
		app.NotFound(w)
		app.ErrorLog.Print(err)
		return
	}

	if product_quantity == 0 {
		app.ErrorMessage(w, 400, "Out of stock")
		return
	}

	if product_quantity < required_quantity {
		app.ErrorMessage(w, 400, "Please, Select less quantity")
		return
	}

	validator = app.User.AddrErrorCheck(input)
	if len(validator) != 0 {
		app.ErrorMessage(w, http.StatusAccepted, validator)
		return
	}

	Order := &models.OrderInfo{}
	Order.AddrId, err = app.User.CreateAddr(input, app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	Order.PaymentMethod = requestData.PaymentMethod
	Order.CartId = requestData.CartID
	Order.ProductId = product_id
	Order.Quantity = required_quantity
	Order.Price = price
	Order.UserId = app.Uid

	order_id, err := app.Order.CreateOrder(Order)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if requestData.PaymentMethod == 1 {
		err = app.Order.ActivateOrder(order_id)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		err = app.Cart.RemoveFromCart(order_id, app.Uid)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// app.sendJSONResponse(w, 200, Order)
		http.Redirect(w, r, fmt.Sprintf("/api/order/review/?order_id=%d", order_id), http.StatusPermanentRedirect)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Now go to Payment page")
}

func (app *Application) MakePayment(w http.ResponseWriter, r *http.Request) {
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	err, active := app.Order.OrderStatus(order_id, app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if active {
		http.Redirect(w, r, fmt.Sprintf("/api/order/review/?order_id=%d", order_id), http.StatusPermanentRedirect)
		return
	}

	quantity, price, err := app.Payment.ValidOrder(app.Uid, order_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = app.Product.ProductQuantity(order_id, quantity)
	if err != nil {
		if err == models.ErrNoRecord {
			app.ErrorMessage(w, 200, "Not enough quantity")
			return
		}

		app.ServerError(w, err)
		return
	}

	total_price := quantity * price
	app.InfoLog.Print(total_price)
	err = app.Order.ActivateOrder(order_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = app.Cart.RemoveFromCart(order_id, app.Uid)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/api/order/review/?order_id=%d", order_id), http.StatusPermanentRedirect)
}

func (app *Application) OrderReview(w http.ResponseWriter, r *http.Request) {
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	if order_id == 0 || err != nil {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	order_info, err := app.Order.OrderInfo(order_id, app.Uid)
	if err != nil {
		if err == models.ErrNoRecord {
			app.NotFound(w)
			app.ErrorLog.Print(err)
			return
		}

		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusFound, &order_info)
}

func (app *Application) UpdateOrderQuantity(w http.ResponseWriter, r *http.Request) {
	order_id, err := strconv.Atoi(r.URL.Query().Get("order_id"))
	quantity, err1 := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err1 != nil || err != nil || order_id == 0 {
		app.NotFound(w)
		app.InfoLog.Print(err)
		return
	}

	err = app.Order.UpdateOrderQuantity(quantity, app.Uid, order_id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.FinalMessage(w, http.StatusAccepted, "Quantity Updated")
}
