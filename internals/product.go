package models

import (
	"database/sql"
	"github.com/iamgak/go-ecommerce/validator"
	"regexp"
)

type ProductModelInterface interface {
	CreateProduct(*Product) error
	CreateProductAddr(*Product_Addr) error
	DeleteProduct(int, int) error
	ChangeProductPrice(float32, int, int) error
	UpdateProduct(*Product, int) error
	UserProductExist(int, int) error
	ProductQuantity(int, int) error
	ProductErrorCheck(*Product) map[string]string
	ProductAddrErrorCheck(*Product_Addr) map[string]string
	UpdateProductQuantity(int, int, int) error
	// ListProductWithParams(int) error
	// ListProduct() error
}

type ProductDB struct {
	DB *sql.DB
}

func (pr *ProductDB) ListProductWithParams(product_id int) error {
	var validId int
	err := pr.DB.QueryRow("SELECT 1 FROM product WHERE id = $1 AND uid = $3", product_id).Scan(&validId)
	if err == sql.ErrNoRows {
		return ErrNoRecord
	}

	return err
}

func (pr *ProductDB) ListProduct() error {
	var validId int
	err := pr.DB.QueryRow("SELECT 1 FROM product WHERE id = $1 AND uid = $3").Scan(&validId)
	if err == sql.ErrNoRows {
		return ErrNoRecord
	}

	return err
}

func (pr *ProductDB) DeleteProduct(user_id, product_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET active = FALSE WHERE id = $1 AND user_id = $2", user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Product Deleted", user_id, product_id)
	}

	return err
}

func (pr *ProductDB) UpdateProductQuantity(quantity, user_id, product_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET quantity = $1 WHERE id = $2 AND user_id = $3", quantity, user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Quantity Changed", user_id, product_id)
	}

	return err
}

func (pr *ProductDB) ChangeProductPrice(price float32, user_id, product_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET price = $1 WHERE id = $2 AND user_id = $3", price, user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Price Changed", user_id, product_id)
	}
	return err
}

func (pr *ProductDB) UpdateProduct(obj *Product, product_id int) error {
	_, err := pr.DB.Exec("UPDATE product SET title = $1,category = $2, sub_category = $3, description = $4, amount= $5, price = $6 WHERE user_id=$7 AND id = $8 AND active = TRUE", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Quantity, obj.Price, obj.Uid, product_id)
	if err == nil {
		err = pr.ProductLog("Quantity Changed", obj.Uid, product_id)
	}
	return err
}

func (pr *ProductDB) ProductErrorCheck(obj *Product) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(obj.Title), "Title", "Title field Cannot be Empty")
	validator.CheckField(validator.NotBlank(obj.Description), "Description", "Description field Cannot be Empty")
	validator.CheckField(obj.Quantity > 0, "Quantity", "Quantity field Cannot be Empty")
	validator.CheckField(obj.Price > 0, "Price", "Price field Cannot be Empty Or 0")
	validator.CheckField(obj.Category > 0, "Category", "Category field Cannot be Empty Or 0")
	validator.CheckField(obj.SubCategory > 0, "SubCategory", "SubCategory field Cannot be Empty Or 0")
	return validator.Errors
}

func (pr *ProductDB) ProductAddrErrorCheck(obj *Product_Addr) map[string]string {
	validator := &validator.Validator{Errors: make(map[string]string)}
	validator.CheckField(validator.NotBlank(obj.Addr), "Addr", "Addr field Cannot be Empty")
	validator.CheckField(validator.NotBlank(obj.Pincode), "Pincode", "Pincode field Cannot be Empty")
	validator.CheckField(validator.NotBlank(obj.Mobile), "Mobile", "Mobile field Cannot be Empty")
	validator.CheckField(obj.Region_id > 0, "Region_id", "Region_id field Cannot be Empty or zero")
	validator.CheckField(obj.District_id > 0, "District_id", "District_id field Cannot be Empty Or 0")
	if len(obj.Pincode) > 0 {
		PincodePattern := regexp.MustCompile(`^[0-9]{5,8}$`)
		validator.CheckField(PincodePattern.MatchString(obj.Pincode), "Pincode", "Pincode Should contain Numeric value only")
	}

	return validator.Errors
}

func (pr *ProductDB) CreateProduct(obj *Product) error {
	_, err := pr.DB.Exec("INSERT INTO product (title,category_id, sub_category_id, quantity, price, user_id, descriptions) VALUES($1,$2,$3,$4,$5,$6,$7)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	return err
}

func (pr *ProductDB) CreateProductAddr(addr *Product_Addr) error {
	_, err := pr.DB.Exec("INSERT INTO product_origin_addr (order_id, mobile, region_id, district_id, pincode, addr) VALUES ($1,$2,$3,$4,$5,$6)", addr.Order_id, addr.Mobile, addr.Region_id, addr.District_id, addr.Pincode, addr.Addr)
	if err == nil {
		_, err = pr.DB.Exec("UPDATE product SET ACTIVE = TRUE WHERE id = $1", addr.Order_id)
	}

	return err
}

func (pr *ProductDB) UserProductExist(product_id, user_id int) error {
	var validId int
	err := pr.DB.QueryRow("SELECT 1 FROM product WHERE id = $1 AND uid = $3", product_id, user_id).Scan(&validId)
	if err == sql.ErrNoRows {
		return ErrNoRecord
	}

	return err
}

func (pr *ProductDB) ProductQuantity(order_id, quantity int) error {
	var validId int
	err := pr.DB.QueryRow("SELECT 1 FROM product INNER JOIN order_listing ON order_listing.product_id = product.id WHERE active = TRUE AND order_listing.id = $1 AND quantity >= $2 ", order_id, quantity).Scan(&validId)
	if err == sql.ErrNoRows {
		return ErrNoRecord
	}

	return err
}

func (pr *ProductDB) ProductLog(activity string, user_id, product_id int) error {
	_, err := pr.DB.Exec("UPDATE product_log SET superseded = TRUE WHERE activity = $1 AND user_id = $2 AND product_id = $3 ", activity, user_id, product_id)
	if err != nil {
		return err
	}

	_, err = pr.DB.Exec("INSERT INTO product_log ( activity,user_id, product_id) VALUES ( $1, $2, $3)", activity, user_id, product_id)
	return err
}
