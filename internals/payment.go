package models

import (
	"database/sql"
	"fmt"
	"github.com/iamgak/go-ecommerce/validator"
)

type PaymentModelInterface interface {
	CreatePayment(*Payment) error
	// CancelOrder(int, int) error
	ErrorCheck(*Payment) map[string]string
}

type Payment struct {
	CartId  int
	Addr    string
	Pincode int
	State   string
	Price   float32
	IpAddr  string
}

type PaymentDB struct {
	DB *sql.DB
}

func (pr *PaymentDB) CreatePayment(payment *Payment) error {
	// _, err := pr.DB.Exec("INSERT INTO payment (product_id,  addr, pincode, state) VALUES ($1,$2,$3,$4)", order.CartId, order.Addr, order.Pincode, order.State)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (pr *PaymentDB) CancelOrder(Uid, orderId int) error {
	_, err := pr.DB.Exec("UPDATE Order SET is_cancelled = 0 WHERE id = $1 AND Uid = $2", Uid, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PaymentDB) OrderStatus(uid, object_id, step int) error {
	var id int
	err := pr.DB.QueryRow("SELECT 1 FROM order_status WHERE step"+fmt.Sprint(step)+" = ? AND object_id = ?, uid = ?", step, object_id, uid).Scan(&id)
	return err
}

func (pr *PaymentDB) ProductExist(product_id, uid int) bool {
	var validId int
	_ = pr.DB.QueryRow("SELECT 1 FROM object WHERE id = $1 AND is_deleted = FALSE AND uid = $2", product_id, uid).Scan(&validId)
	return validId > 0
}

// ValidPayment checks if a product with the given product_id has a valid payment associated with it.
func (pr *PaymentDB) ValidPayment(productID, userID int) (error, bool) {
	// Prepare a variable to store the result of the query
	var validPaymentExists int

	// SQL query to check if there's a valid payment for the specified product_id
	query := `
		SELECT 1 
		FROM product 
		INNER JOIN user_listing ON user_listing.id = product.uid 
		INNER JOIN cart ON cart.product_id = product.id 
		INNER JOIN order_listing ON order_listing.cart_id = cart.id 
		WHERE product.id = $1 AND product.uid = $2 AND `

	// Execute the query and scan the result into validPaymentExists
	err := pr.DB.QueryRow(query, productID, userID).Scan(&validPaymentExists)
	return err, validPaymentExists > 0
}

func (pr *PaymentDB) CartExist(product_id int) bool {
	var validId int
	_ = pr.DB.QueryRow("SELECT object_id FROM cart WHERE id = $1 AND completed =  ", product_id).Scan(&validId)
	return validId > 0
}

func (pr *PaymentDB) ErrorCheck(payment *Payment) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(payment.Addr), "Address", "Address field Cannot be Empty")
	validator.CheckField(validator.NotBlank(payment.State), "State", "State field Cannot be Empty")
	validator.CheckField(payment.Pincode > 0, "Pincode", "Pincode field Cannot be zero")
	return validator.Errors
}
