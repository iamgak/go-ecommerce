package order

import "database/sql"

type OrderInterf interface {
	CreateOrder(*Order) error
	CancelOrder(int, int) error
}

type Order struct {
	ProductId int
	Uid       int
	Quantity  int
}

type OrderDB struct {
	DB *sql.DB
}

func (pr *OrderDB) CreateOrder(order *Order) error {
	_, err := pr.DB.Exec("INSERT INTO `Order`(`product_id`,`uid`, `quantity`) VALUES($1,$2,$3)", order.ProductId, order.Uid, order.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (pr *OrderDB) CancelOrder(Uid, orderId int) error {
	_, err := pr.DB.Exec("UPDATE `Order` SET `is_cancelled` = 0 WHERE `id` = $1 AND `Uid` = $2", Uid, orderId)
	if err != nil {
		return err
	}

	return nil
}
