package object

import (
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
)

type ObjectInterface interface {
	InsertObject(*Object) error
	UpdateObject(*Object, int) error
	DeleteObject(int, int) error
	ChangePriceObject(float32, int, int) error
}

type Object struct {
	Title       string
	Category    int
	SubCategory int
	Uid         int
	Quantity    int
	Price       float32
	Description string
}

type ObjectDB struct {
	DB *sql.DB
}

func (pr *ObjectDB) ListObject(obj *Object) error {
	_, err := pr.DB.Exec("INSERT INTO `object`(`title`, `quantity`, `category`, `sub-category`, `descriptions`, `price`, `uid`, `last_updated`) VALUES ($1,$2,$3,$4,$5,$6,$7,CURRENT_TIMESTAMP())", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Uid, obj.Price, obj.Price)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) DeleteObject(uid, object_id int) error {
	_, err := pr.DB.Exec("UPDATE `object` SET `is_deleted` = 0 WHERE `id` = $1 AND `Uid` = $2", uid, object_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) ChangePriceObject(price float32, uid, object_id int) error {
	_, err := pr.DB.Exec("UPDATE `object` SET `price` = $1 WHERE `id` = $2 AND `Uid` = $3", price, uid, object_id)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ObjectDB) UpdateObject(obj *Object, object_id int) error {
	_, err := pr.DB.Exec("UPDATE `object` SET `title` = $1,`Category`= $2, `SubCategory`= $3, `Description`= $4, `amount`= $5, `Cost` = $6 WHERE Uid=$7 AND id = $8", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Quantity, obj.Price, obj.Uid, object_id)
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
	_, err := pr.DB.Exec("INSERT INTO object (title,category,sub_category,quantity, price, uid,descriptions) VALUES($1,$2,$3,$4,$5,$6,$7)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	return err
}
