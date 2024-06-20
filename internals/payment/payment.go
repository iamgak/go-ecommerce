package payment

import (
	"database/sql"
	"fmt"

	"github.com/iamgak/go-ecommerce/validator"
)

type PaymentFunc interface {
	CreateOrder(*Order) error
	CancelOrder(int, int) error
}

type Order struct {
	CartId  int
	Addr    string
	Pincode int
	State   string
	Price   float32
	IpAddr  string
}

type OrderDB struct {
	DB *sql.DB
}

func (pr *OrderDB) CreateOrder(order *Order) error {
	_, err := pr.DB.Exec("INSERT INTO order (product_id,  addr, pincode, state) VALUES ($1,$2,$3,$4)", order.CartId, order.Addr, order.Pincode, order.State)
	if err != nil {
		return err
	}

	return nil
}

func (pr *OrderDB) CancelOrder(Uid, orderId int) error {
	_, err := pr.DB.Exec("UPDATE Order SET is_cancelled = 0 WHERE id = $1 AND Uid = $2", Uid, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (pr *OrderDB) OrderStatus(uid, object_id, step int) error {
	var id int
	err := pr.DB.QueryRow("SELECT 1 FROM order_status WHERE step"+fmt.Sprint(step)+" = ? AND object_id = ?, uid = ?", step, object_id, uid).Scan(&id)
	return err
}

func (pr *OrderDB) ProductExist(product_id, uid int) bool {
	var validId int
	_ = pr.DB.QueryRow("SELECT 1 FROM object WHERE id = $1 AND is_deleted = FALSE AND uid = $2", product_id, uid).Scan(&validId)
	return validId > 0
}

func (pr *OrderDB) OrderExist(product_id int) bool {
	var validId int
	_ = pr.DB.QueryRow("SELECT object_id FROM order WHERE id = $1 AND completed =  ", product_id).Scan(&validId)
	return validId > 0
}

func (pr *OrderDB) CartExist(product_id int) bool {
	var validId int
	_ = pr.DB.QueryRow("SELECT object_id FROM cart WHERE id = $1 AND completed =  ", product_id).Scan(&validId)
	return validId > 0
}

func (pr *OrderDB) ErrorCheck(order *Order) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(order.Addr), "Address", "Address field Cannot be Empty")
	validator.CheckField(validator.NotBlank(order.State), "State", "State field Cannot be Empty")
	validator.CheckField(order.Pincode > 0, "Pincode", "Pincode field Cannot be zero")
	if len(validator.Errors) == 0 {
		// validator.CheckField(pr.CartExist(order.CartId), "Warning", "Invalid Request")
	}

	return validator.Errors
}
