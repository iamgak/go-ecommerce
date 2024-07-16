package data

import (
	"context"
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
)

type OrderDB struct {
	db  *sql.DB
	ctx context.Context
}

func (c *OrderDB) Close() {
	c.db.Close()
}

func (c *OrderDB) OrderListing(user_id int) ([]*OrderReview, error) {
	query := `SELECT type_payment.name, order_listing.is_cancelled, 
				order_listing.price, order_listing.active, order_listing.dispatch,
				product.title, order_listing.created_at
				FROM order_listing
				INNER JOIN product ON product.id = order_listing.product_id
				INNER JOIN type_payment ON type_payment.id = order_listing.payment_method
				WHERE order_listing.user_id = $1`

	return c.Listing(query)
}

func (ord *OrderDB) Listing(stmt string) ([]*OrderReview, error) {

	rows, err := ord.db.QueryContext(ord.ctx, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	defer rows.Close()
	defer ord.Close()
	arr := []*OrderReview{}
	for rows.Next() {
		bk, err := ord.ScanData(rows)
		if err != nil {
			return nil, err
		}

		arr = append(arr, bk)
	}

	return arr, err
}

func (m *OrderDB) ScanData(rows *sql.Rows) (*OrderReview, error) {
	info := &OrderReview{}
	arg := []interface{}{
		&info.PaymentType,
		&info.Cancelled,
		&info.ProductPrice,
		&info.Active,
		&info.Dispatched,
		&info.ProductName,
		&info.OrderAt,
	}

	err := rows.Scan(arg...)
	return info, err
}

func (ord *OrderDB) CreateOrder(order *OrderInfo) (int, error) {
	var order_id int
	err := ord.db.QueryRowContext(ord.ctx, "INSERT INTO order_listing (cart_id,product_id, quantity, price, addr_id,payment_method, user_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id", order.CartId, order.ProductId, order.Quantity, order.Price, order.AddrId, order.PaymentMethod, order.UserId).Scan(&order_id)
	return order_id, err
}

func (ord *OrderDB) CancelOrder(Uid, orderId int) error {

	err := ord.db.QueryRowContext(ord.ctx, "UPDATE order_listing SET is_cancelled = TRUE WHERE id = $1 AND user_id = $2 RETURNING id", orderId, Uid).Scan(&orderId)
	if err != nil {
		return err
	}

	_, err = ord.db.Exec("INSERT INTO order_cancel (order_id) VALUES ($1)", orderId)
	return err
}

func (ord *OrderDB) UpdateOrderQuantity(quantity, user_id, orderId int) error {
	_, err := ord.db.Exec("UPDATE order_listing SET quantity = $1 WHERE id = $2 AND user_id = $3", quantity, orderId, user_id)
	return err
}

func (ord *OrderDB) ValidCart(cartId, user_id int) (int, int, int, float32, error) {
	var product_id, product_quantity, cart_quantity int
	var price float32
	query := `SELECT product.id, product.quantity, cart.quantity, product.price  
				FROM product 
				INNER JOIN cart ON cart.product_id = ordoduct.id 
				WHERE cart.user_id = $1 AND cart.id = $2 AND product.active = TRUE`

	arg := []interface{}{
		&product_id,
		&product_quantity,
		&cart_quantity,
		&price,
	}

	err := ord.db.QueryRowContext(ord.ctx, query, user_id, cartId).Scan(arg...)
	if err != nil && err == sql.ErrNoRows {
		return product_id, product_quantity, cart_quantity, price, ErrNoRecord
	}

	return product_id, product_quantity, cart_quantity, price, err
}

func (c *OrderDB) OrderInfo(orderId, user_id int) ([]*OrderReview, error) {
	query := `SELECT type_payment.name, order_listing.is_cancelled, 
				order_listing.price, order_listing.active, order_listing.dispatch,
				product.title, order_listing.created_at
				FROM order_listing
				INNER JOIN product ON product.id = order_listing.product_id
				INNER JOIN type_payment ON type_payment.id = order_listing.payment_method
				WHERE order_listing.user_id = $1 AND order_listing.id = $2`
	return c.Listing(query)
}

// cases whern then in postgres
func (ord *OrderDB) OrderStatus(uid, product_id int) (error, bool) {
	var active bool

	err := ord.db.QueryRowContext(ord.ctx, "SELECT active FROM order_listing WHERE id= $1 AND user_id = $2", product_id, uid).Scan(&active)
	return err, active
}

func (ord *OrderDB) ActivateOrder(order_id int) error {
	var quantity, product_id, cart_id int

	err := ord.db.QueryRowContext(ord.ctx, "UPDATE order_listing SET active = TRUE WHERE id = $1 RETURNING quantity, product_id, cart_id", order_id).Scan(&quantity, &product_id, &cart_id)
	if err == nil {
		_, err = ord.db.Exec("UPDATE product SET quantity = quantity-$1 WHERE id = $2", quantity, product_id)
		if err != nil {
			return err
		}

		_, err = ord.db.Exec("UPDATE cart SET active = FALSE WHERE id = $1", cart_id)
	}

	return err
}

func (ord *OrderDB) RequestErrorCheck(data *RequestData) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(data.CartID > 0, "cart_id", "cart_id field Cannot be Empty")
	validator.CheckField(data.PaymentMethod > 0, "payment_method", "payment_method field is Empty")
	if len(validator.Errors) == 0 {
		validator.CheckField(ord.CartExist(data.CartID, data.UserId), "cart_id", "Invalid cart_id")
	}

	if len(validator.Errors) == 0 {
		validator.CheckField(ord.PaymentMethodExist(data.PaymentMethod), "payment_method", "Invalid payment_method")
	}

	if data.AddrId > 0 && len(validator.Errors) > 0 {
		validator.CheckField(ord.UserAddrExist(data.AddrId, data.UserId), "addr_id", "Invalid address id")
	}

	return validator.Errors
}

func (ord *OrderDB) CartExist(cart_id, user_id int) bool {
	var validId int

	_ = ord.db.QueryRowContext(ord.ctx, "SELECT 1 FROM cart WHERE id = $1 AND user_id =$2 AND active = TRUE", cart_id, user_id).Scan(&validId)
	return validId > 0
}

func (ord *OrderDB) PaymentMethodExist(payment_id int) bool {
	var validId int

	_ = ord.db.QueryRowContext(ord.ctx, "SELECT 1 FROM type_payment WHERE id = $1", payment_id).Scan(&validId)
	return validId > 0
}

func (ord *OrderDB) UserAddrExist(address_id, user_id int) bool {
	var validId int

	_ = ord.db.QueryRowContext(ord.ctx, "SELECT 1 FROM user_addr WHERE id = $1 AND user_id = $2", address_id, user_id).Scan(&validId)
	return validId > 0
}
