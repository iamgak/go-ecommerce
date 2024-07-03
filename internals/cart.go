package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/iamgak/go-ecommerce/validator"
)

type CartModelInterface interface {
	AddInCart(*Cart) error
	RemoveFromCart(int, int) error
	ErrorCheck(*Cart) map[string]string
	ProductExist(int) bool
	CartListing(int) ([]*Listing, error)
}

type CartDB struct {
	DB *sql.DB
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

func (m *CartDB) Listing(query string, user_id int) ([]*Listing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, user_id)
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
		bk, err := m.ScanData(rows)
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

func (m *CartDB) ScanData(rows *sql.Rows) (*Listing, error) {
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
	_, err := c.DB.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1,$2,$3) ", cart.Uid, cart.ProductId, cart.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartDB) RemoveFromCart(CartId, user_id int) error {
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
