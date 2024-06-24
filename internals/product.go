package models

import (
	"database/sql"
	"log"

	"github.com/iamgak/go-ecommerce/validator"
)

type ProductModelInterface interface {
	CreateObject(*Object) error
	DeleteObject(int, int) error
	ChangePriceObject(float32, int, int) error
	UpdateObject(*Object, int) error
	ListObject(*Object) error
	ErrorCheck(*Object) map[string]string
}

type Object struct {
	Title       string
	Category    int
	SubCategory int
	Uid         int
	Quantity    int
	Addr_id     int
	Price       float32
	Description string
}

type ObjectDB struct {
	DB *sql.DB
}

func (pr *ObjectDB) ListObject(obj *Object) error {
	_, err := pr.DB.Exec("INSERT INTO product (title, quantity, category_id, sub_category_id, descriptions, price, user_id, origin_addr_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Uid, obj.Price, obj.Price, obj.Addr_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) DeleteObject(user_id, object_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET is_deleted = 0 WHERE id = $1 AND user_id = $2", user_id, object_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) ChangePriceObject(price float32, user_id, object_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET price = $1 WHERE id = $2 AND user_id = $3", price, user_id, object_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) UpdateObject(obj *Object, object_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET title = $1,category = $2, sub_category = $3, description = $4, amount= $5, price = $6 WHERE user_id=$7 AND id = $8", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Quantity, obj.Price, obj.Uid, object_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) ErrorCheck(obj *Object) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(obj.Title), "Title", "Title field Cannot be Empty")
	validator.CheckField(validator.NotBlank(obj.Description), "Description", "Description field Cannot be Empty")
	validator.CheckField(obj.Quantity > 0, "Amount", "Amount field Cannot be Empty")
	validator.CheckField(obj.Price > 0, "Price", "Price field Cannot be Empty Or 0")
	validator.CheckField(obj.Category > 0, "Category", "Category field Cannot be Empty Or 0")
	validator.CheckField(obj.SubCategory > 0, "SubCategory", "SubCategory field Cannot be Empty Or 0")
	return validator.Errors
}

func (pr *ObjectDB) CreateObject(obj *Object) error {
	_, err := pr.DB.Exec("INSERT INTO product (title,category,sub_category,quantity, price, user_id,descriptions) VALUES($1,$2,$3,$4,$5,$6,$7)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	log.Printf("INSERT INTO product (title,category,sub_category,quantity, price, user_id,descriptions) VALUES(%s,%d,%d,%d,%f,%d,%s)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	return err
}
