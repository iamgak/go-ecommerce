package models

import (
	"database/sql"

	"github.com/iamgak/go-ecommerce/validator"
)

type CartModelInterface interface {
	AddInCart(*Cart) error
	DeleteItem(int, int) error
	ErrorCheck(*Cart) map[string]string
	ProductExist(int) bool
}
type Cart struct {
	ProductId int
	Uid       int
	Quantity  int
	Status    bool
}

type CartDB struct {
	DB *sql.DB
}

func (c *CartDB) AddInCart(cart *Cart) error {
	_, err := c.DB.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1,$2,$3) ", cart.Uid, cart.ProductId, cart.Quantity)
	if err != nil {
		return ErrCantAddInCart
	}

	return nil
}

func (c *CartDB) DeleteItem(Uid, CartId int) error {
	_, err := c.DB.Exec("DELETE cart WHERE id = $1 AND user_id = $2", CartId, Uid)
	if err != nil {
		return ErrCantRemoveItem
	}

	return nil
}

func (c *CartDB) ErrorCheck(cart *Cart) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(cart.Quantity > 0, "Quantity", "Quantity field Cannot be Empty")
	validator.CheckField(cart.ProductId > 0, "ProductId", "ProductId field is Empty")
	if len(validator.Errors) == 0 {
		validator.CheckField(c.ProductExist(cart.ProductId), "Warning", "Invalid Request")
	}

	return validator.Errors
}

func (c *CartDB) ProductExist(product_id int) bool {
	var validId int
	_ = c.DB.QueryRow("SELECT 1 FROM product WHERE id = $1 AND is_deleted = FALSE", product_id).Scan(&validId)
	return validId > 0
}
