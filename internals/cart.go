package products

import "database/sql"

type Cart struct {
	ProductId int
	Uid       int
	Title     string
	Status    bool
}

type CartDB struct {
	DB *sql.DB
}

func (cart *CartDB) InsertItem(product *Products) error {
	_, err := cart.DB.Exec("INSERT INTO `Products`(`title`,`Category`, `SubCategory`, `Description`, `Uid`, `Instock`,`Cost`) VALUES($1,$2,$3,$4,$5,$6,$7)", product.Title, product.Category, product.SubCategory, product.Description, product.Uid, product.InStock, product.Cost)
	if err != nil {
		return err
	}

	return nil
}

func (cart *CartDB) DeleteItem(Uid, CartId int) error {
	_, err := cart.DB.Exec("DELETE `Cart` WHERE `id` = $1 AND `Uid` = $2", Uid, CartId)
	if err != nil {
		return err
	}

	return nil
}
