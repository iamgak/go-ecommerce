package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"sync"
	"time"

	"github.com/iamgak/go-ecommerce/validator"
	"github.com/redis/go-redis/v9"
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
	ProductListing() ([]*ProductListing, error)

	// ListProduct() error
}

type ProductDB struct {
	db    *sql.DB
	redis *redis.Client
	mute  *sync.RWMutex
	ctx   context.Context
}

func (c *ProductDB) Close() {
	c.redis.Close()
	c.db.Close()
}
func (c *ProductDB) ProductListing() ([]*ProductListing, error) {
	c.mute.RLock()
	defer c.mute.Unlock()
	query := `SELECT pt.title, ct.category, pt.quantity, 
				pt.price, pt.active 
				FROM product pt
				INNER JOIN category_main ct ON pt.category_id = ct.id
				WHERE pt.active = TRUE`

	cart_listing, err := c.Listing(query)
	return cart_listing, err
}

func (m *ProductDB) Listing(stmt string) ([]*ProductListing, error) {

	products := []*ProductListing{}
	queryBytes, err := json.Marshal(stmt)
	if err != nil {
		panic(err)
	}

	val, err := m.redis.Get(m.ctx, string(queryBytes)).Result()
	if err == nil {
		// Deserialize the cached result
		err = json.Unmarshal([]byte(val), &products)
		if err != nil {
			return nil, err
		}

		return products, err
	} else if err != redis.Nil {
		return nil, err
	}

	rows, err := m.db.Query(stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		arr := new(ProductListing)
		arr, err := m.ScanData(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, arr)
	}

	// Cache the result in Redis for 5 minutes
	data, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	err = m.redis.Set(m.ctx, string(queryBytes), data, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	m.Close()
	return products, err

}

func (m *ProductDB) ScanData(rows *sql.Rows) (*ProductListing, error) {
	listing := new(ProductListing)
	arg := []interface{}{
		&listing.Title,
		&listing.Category,
		&listing.Quantity,
		&listing.Price,
		&listing.Available,
	}

	err := rows.Scan(arg...)
	return listing, err
}

func (pr *ProductDB) DeleteProduct(user_id, product_id int) error {
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("UPDATE product SET active = FALSE WHERE id = $1 AND user_id = $2", user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Product Deleted", user_id, product_id)
	}

	return err
}

func (pr *ProductDB) UpdateProductQuantity(quantity, user_id, product_id int) error {
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("UPDATE product SET quantity = $1 WHERE id = $2 AND user_id = $3", quantity, user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Quantity Changed", user_id, product_id)
	}

	return err
}

func (pr *ProductDB) ChangeProductPrice(price float32, user_id, product_id int) error {
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("UPDATE product SET price = $1 WHERE id = $2 AND user_id = $3", price, user_id, product_id)
	if err == nil {
		err = pr.ProductLog("Price Changed", user_id, product_id)
	}
	return err
}

func (pr *ProductDB) UpdateProduct(obj *Product, product_id int) error {
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("UPDATE product SET title = $1,category = $2, sub_category = $3, description = $4, amount= $5, price = $6 WHERE user_id=$7 AND id = $8 AND active = TRUE", obj.Title, obj.Category, obj.SubCategory, obj.Description, obj.Quantity, obj.Price, obj.Uid, product_id)
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
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("INSERT INTO product (title,category_id, sub_category_id, quantity, price, user_id, descriptions) VALUES($1,$2,$3,$4,$5,$6,$7)", obj.Title, obj.Category, obj.SubCategory, obj.Quantity, obj.Price, obj.Uid, obj.Description)
	return err
}

func (pr *ProductDB) CreateProductAddr(addr *Product_Addr) error {
	pr.mute.Lock()
	defer pr.mute.Unlock()
	_, err := pr.db.Exec("INSERT INTO product_origin_addr (order_id, mobile, region_id, district_id, pincode, addr) VALUES ($1,$2,$3,$4,$5,$6)", addr.Order_id, addr.Mobile, addr.Region_id, addr.District_id, addr.Pincode, addr.Addr)
	if err == nil {
		_, err = pr.db.Exec("UPDATE product SET ACTIVE = TRUE WHERE id = $1", addr.Order_id)
	}

	return err
}

func (pr *ProductDB) UserProductExist(product_id, user_id int) error {
	var validId int
	err := pr.db.QueryRowContext(pr.ctx, "SELECT 1 FROM product WHERE id = $1 AND uid = $3", product_id, user_id).Scan(&validId)
	if err == sql.ErrNoRows {
		return ErrNoRecord
	}

	return err
}

func (pr *ProductDB) ProductQuantity(order_id, quantity int) error {
	var validId int
	err := pr.db.QueryRowContext(pr.ctx, "SELECT 1 FROM product INNER JOIN order_listing ON order_listing.product_id = product.id WHERE active = TRUE AND order_listing.id = $1 AND quantity >= $2 ", order_id, quantity).Scan(&validId)
	return err
}

func (pr *ProductDB) ProductLog(activity string, user_id, product_id int) error {
	_, err := pr.db.Exec("UPDATE product_log SET superseded = TRUE WHERE activity = $1 AND user_id = $2 AND product_id = $3 ", activity, user_id, product_id)
	if err != nil {
		return err
	}

	_, err = pr.db.Exec("INSERT INTO product_log ( activity,user_id, product_id) VALUES ( $1, $2, $3)", activity, user_id, product_id)
	return err
}
