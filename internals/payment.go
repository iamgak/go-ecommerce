package data

import (
	"database/sql"
)

type PaymentModelInterface interface {
	CreatePayment(*Payment) error
	ValidOrder(int, int) (int, int, error)
	// ErrorCheck(*Payment) map[string]string
}

type PaymentDB struct {
	DB *sql.DB
}

func (pr *PaymentDB) CreatePayment(payment *Payment) error {
	var transaction_id int
	err := pr.DB.QueryRow("INSERT INTO order_payment (order_id, transaction_id, amount, status) VALUES ($1,$2,$3,$4) RETURNING id", payment.OrderId, payment.TransactionId, payment.Amount, payment.Status).Scan(&transaction_id)
	// if err == nil {
	// 	_, err = pr.DB.Exec("UPDATE order_listing SET transaction_id = $1 WHERE id = $2 ", transaction_id, payment.OrderId)
	// }

	return err
}

func (pr *PaymentDB) ValidOrder(userId, OrderId int) (int, int, error) {
	var quantity, price int
	query := `SELECT price, quantity
				FROM product 
				INNER JOIN order_listing ON order_listing.product_id = product.id 
				WHERE order.id = $1 AND order_listing.uid = $2 AND  active = FALSE AND order.quantity <= product.quantity`

	err := pr.DB.QueryRow(query, OrderId, userId).Scan(&price, &quantity)
	if err != nil && err == sql.ErrNoRows {
		return quantity, price, ErrNoRecord
	}

	return quantity, price, err
}
