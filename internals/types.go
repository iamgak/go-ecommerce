package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Seller struct {
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Password    string `json:"password"`
	Pancard     string `json:"pancard"`
	Mobile      string `json:"mobile"`
	Region_id   int    `json:"region"`
	District_id int    `json:"district"`
	Pincode     string `json:"pincode"`
	Addr        string `json:"addr"`
}

type ForgetPassword struct {
	Email string `json:"email"`
}

type NewPassword struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UserAddr struct {
	Pincode     string `json:"pincode"`
	Region_id   int    `json:"region"`
	District_id int    `json:"district"`
	Addr        string `json:"addr"`
	Mobile      string `json:"mobile"`
}

type Product struct {
	Title       string  `json:"title"`
	Category    int     `json:"category"`
	SubCategory int     `json:"sub_category"`
	Uid         int     `json:"uid"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
}

type Product_Addr struct {
	Order_id    int    `json:"order_id"`
	Region_id   int    `json:"region"`
	District_id int    `json:"district"`
	Pincode     string `json:"pincode"`
	Mobile      string `json:"mobile"`
	Addr        string `json:"addr"`
}

type OrderInfo struct {
	ProductId     int     `json:"product_id"`
	CartId        int     `json:"cart_id"`
	Quantity      int     `json:"quantity"`
	Price         float32 `json:"price"`
	AddrId        int     `json:"addr_id"`
	UserId        int     `json:"user_id"`
	PaymentMethod int     `json:"payment_method"`
}

type Cart struct {
	ProductId int `json:"product_id"`
	Uid       int `json:"_"`
	Quantity  int `json:"quantity"`
}

type Payment struct {
	OrderId       int
	Amount        float32
	TransactionId string
	Status        bool
}

type OrderReview struct {
	PaymentId    int
	ProductName  string
	ProductPrice float32
	PaymentType  string
	Active       bool
	Dispatched   bool
	Cancelled    bool
	OrderAt      string
}

type RequestData struct {
	CartID        int `json:"cart_id"`
	PaymentMethod int `json:"payment_method"`
	AddrId        int `json:"addr_id"`
	UserId        int `json:"user_id"`
}
