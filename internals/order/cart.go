package order

import "database/sql"

type Cart struct {
	ProductId int
	Uid       int
	Title     string
	Price     float32
	Status    bool
}

type CartDB struct {
	DB *sql.DB
}

func (c *CartDB) InsertItem(cart *Cart) error {
	_, err := c.DB.Exec("INSERT INTO `UserCartListing`(`uid`,`pid`, `price`) VALUES($1,$2,$3)", cart.Uid, cart.ProductId, cart.Price)
	if err != nil {
		return err
	}

	return nil
}

func (cart *CartDB) DeleteItem(Uid, CartId int) error {
	_, err := cart.DB.Exec("DELETE `UserCartListing` WHERE `id` = $1 AND `Uid` = $2", Uid, CartId)
	if err != nil {
		return err
	}

	return nil
}
