package data

import (
	"context"
	"database/sql"
	"sync"

	"github.com/iamgak/go-ecommerce/validator"
)

type CartDB struct {
	db   *sql.DB
	ctx  context.Context
	mute sync.RWMutex
}

func (c *CartDB) CartListing(user_id int) ([]*Listing, error) {
	query := `SELECT pt.title, ct.id, pt.quantity, 
				cart.quantity, pt.price, pt.active 
				FROM product pt
				INNER JOIN category_main ct ON pt.category_id = ct.id
				INNER JOIN cart ON pt.id = cart.product_id
				WHERE cart.user_id = $1 AND cart.active = TRUE`

	cart_listing, err := c.Listing(query, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return cart_listing, nil
}

func (c *CartDB) Listing(query string, user_id int) ([]*Listing, error) {
	c.mute.RLock()
	defer c.mute.Unlock()
	rows, err := c.db.QueryContext(c.ctx, query, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	Books := []*Listing{}
	for rows.Next() {
		bk, err := c.ScanData(rows)
		if err != nil {
			return nil, err
		}

		Books = append(Books, bk)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Books, err
}

func (c *CartDB) ScanData(rows *sql.Rows) (*Listing, error) {
	listing := new(Listing)
	arg := []interface{}{
		listing.Title,
		&listing.Category,
		&listing.Quantity,
		&listing.SelectedQuantity,
		&listing.Price,
		&listing.Available,
	}
	err := rows.Scan(arg...)
	return listing, err
}

func (c *CartDB) AddInCart(cart *Cart) error {
	c.mute.Lock()
	defer c.mute.Unlock()
	_, err := c.db.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1,$2,$3) ", cart.Uid, cart.ProductId, cart.Quantity)
	return err
}

func (c *CartDB) RemoveFromCart(CartId, user_id int) error {
	c.mute.Lock()
	defer c.mute.Unlock()
	_, err := c.db.Exec("UPDATE cart SET active = FALSE WHERE id = $1 AND user_id = $2", CartId, user_id)
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
	c.mute.RLock()
	defer c.mute.Unlock()
	_ = c.db.QueryRowContext(c.ctx, "SELECT 1 FROM product WHERE id = $1 AND active = TRUE", product_id).Scan(&validId)
	return validId > 0
}
