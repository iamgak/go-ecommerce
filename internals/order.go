package models

import (
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
)

type OrderModelInterface interface {
	CreateOrder(*Order) error
	CancelOrder(int, int) error
	ErrorCheck(*Order) map[string]string
}

type Order struct {
	CartId     int
	Addr       string
	Pincode    int
	State      string
	District   string
	TotalPrice int
	Mobile     string
}

type OrderDB struct {
	DB *sql.DB
}

func (pr *OrderDB) CreateOrder(order *Order) error {
	var order_id int
	err := pr.DB.QueryRow("SELECT id FROM order_listing WHERE cart_id = $1", order.CartId).Scan(&order_id)
	if err != nil {
		if err != sql.ErrNoRows {
			return ErrCantGetItem
		}
		return err
	}

	if order_id > 0 {
		_, err = pr.DB.Exec("UPDATE order_addr SET addr = $2, region = $3, district = $4, pincode = $5,mobile = $6 WHERE order_id = $1", order_id, order.Addr, order.State, order.District, order.Pincode, order.Mobile)
	} else {
		err = pr.DB.QueryRow("INSERT INTO order_listing (cart_id) VALUES ($1) RETURNING id", order.CartId).Scan(&order_id)
		if err != nil {
			return err
		}

		_, err = pr.DB.Exec("INSERT INTO order_addr (order_id, addr, region, district, pincode,mobile) VALUES ($1,$2, $3, $4, $5, $6)", order_id, order.Addr, order.State, order.District, order.Pincode, order.Mobile)
	}

	return err
}

func (pr *OrderDB) CancelOrder(Uid, orderId int) error {
	_, err := pr.DB.Exec("UPDATE order_listing SET is_cancelled = 0 WHERE id = $1 AND user_id = $2", Uid, orderId)
	return err
}

func (pr *OrderDB) ErrorCheck(order *Order) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(order.Addr), "Address", "Address field Cannot be Empty")
	validator.CheckField(validator.NotBlank(order.State), "State", "State field Cannot be Empty")
	validator.CheckField(validator.NotBlank(order.District), "District", "District field Cannot be Empty")
	validator.CheckField(validator.NotBlank(order.Mobile), "Mobile", "Mobile field Cannot be Empty")
	validator.CheckField(order.Pincode > 0, "Pincode", "Pincode field Cannot be zero")
	if len(validator.Errors) == 0 {
		validator.CheckField(validator.ValidNumber(order.Mobile), "Mobile", "Invalid Mobile No format")
	}

	return validator.Errors
}
