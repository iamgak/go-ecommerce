package cart

import (
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
)

type Cart struct {
	ProductId int
	Uid       int
	Price     float32
	Quantity  int
	Status    bool
}

type CartDB struct {
	DB *sql.DB
}

func (c *CartDB) AddInCart(cart *Cart) error {
	_, err := c.DB.Exec("INSERT INTO cart (uid, object_id, price, quantity) VALUES ($1,$2,$3,$4) ", cart.Uid, cart.ProductId, cart.Price, cart.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartDB) DeleteItem(Uid, CartId int) error {
	_, err := c.DB.Exec("DELETE cart WHERE id = $1 AND uid = $2", CartId, Uid)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartDB) ErrorCheck(cart *Cart) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(cart.Quantity > 0, "Quantity", "Quantity field Cannot be Empty")
	validator.CheckField(cart.ProductId > 0, "ProductId", "ProductId field is Empty")
	validator.CheckField(cart.Price > 0, "Price", "Price field Cannot be Empty Or 0")
	if len(validator.Errors) == 0 {
		validator.CheckField(c.ProductExist(cart.ProductId), "Warning", "Invalid Request")
	}

	return validator.Errors
}

func (c *CartDB) CreateObject(cart *Cart) error {
	// _, err := c.DB.Exec("INSERT INTO object (title,category,sub_category,quantity, price, uid,descriptions) VALUES($1,$2,$3,$4,$5,$6,$7)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	// log.Printf("INSERT INTO object (title,category,sub_category,quantity, price, uid,descriptions) VALUES(%s,%d,%d,%d,%f,%d,%s)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	return nil
}

func (c *CartDB) ProductExist(product_id int) bool {
	var validId int
	_ = c.DB.QueryRow("SELECT 1 FROM object WHERE id = $1 AND is_deleted = FALSE ", product_id).Scan(&validId)
	return validId > 0
}
