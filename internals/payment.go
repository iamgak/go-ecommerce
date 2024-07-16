package data

import (
	"context"
	"database/sql"
)

type PaymentModelInterface interface {
	CreatePayment(*Payment) error
	ValidOrder(int, int) (int, int, error)
}

type PaymentDB struct {
	db  *sql.DB
	ctx context.Context
}

func (pr *PaymentDB) CreatePayment(payment *Payment) error {
	var transaction_id int
	err := pr.db.QueryRow("INSERT INTO order_payment (order_id, transaction_id, amount, status) VALUES ($1,$2,$3,$4) RETURNING id", payment.OrderId, payment.TransactionId, payment.Amount, payment.Status).Scan(&transaction_id)
	return err
}

func (pr *PaymentDB) ValidOrder(userId, OrderId int) (int, int, error) {
	var quantity, price int
	query := `SELECT price, quantity
				FROM product 
				INNER JOIN order_listing ON order_listing.product_id = product.id 
				WHERE order.id = $1 AND order_listing.uid = $2 AND  active = FALSE AND order.quantity <= product.quantity`

	err := pr.db.QueryRow(query, OrderId, userId).Scan(&price, &quantity)
	if err != nil && err == sql.ErrNoRows {
		return quantity, price, ErrNoRecord
	}

	return quantity, price, err
}
