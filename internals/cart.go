package models

import (
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
)

type CartModelInterface interface {
	AddInCart(*Cart) error
	RemoveFromCart(int, int) error
	ErrorCheck(*Cart) map[string]string
	ProductExist(int) bool
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

func (c *CartDB) RemoveFromCart(CartId, user_id int) error {
	// inner join in update postgres
	_, err := c.DB.Exec("UPDATE cart SET active = FALSE WHERE id = $1 AND user_id = $2", CartId, user_id)
	return err
}

func (c *CartDB) ErrorCheck(cart *Cart) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(cart.Quantity > 0, "quantity", "quantity field Cannot be Empty")
	validator.CheckField(cart.ProductId > 0, "product_id", "product_id field is Empty")
	if len(validator.Errors) == 0 {
		validator.CheckField(c.ProductExist(cart.ProductId), "Product_id", "Invalid Product_id")
	}

	return validator.Errors
}

func (c *CartDB) ProductExist(product_id int) bool {
	var validId int
	_ = c.DB.QueryRow("SELECT 1 FROM product WHERE id = $1 AND active = TRUE", product_id).Scan(&validId)
	return validId > 0
}
