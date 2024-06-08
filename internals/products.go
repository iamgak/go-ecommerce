package products

import "database/sql"

type Product interface {
	InsertItem(*Products) error
	UpdateItem(*Products, int) error
	DeleteItem(int, int) error
}

type Products struct {
	Title       string
	Category    string
	SubCategory string
	Description string
	Uid         int
	InStock     int
	Cost        float32
}

type ProductDB struct {
	DB *sql.DB
}

func (pr *ProductDB) InsertItem(product *Products) error {
	_, err := pr.DB.Exec("INSERT INTO `Products`(`title`,`Category`, `SubCategory`, `Description`, `Uid`, `Instock`,`Cost`) VALUES($1,$2,$3,$4,$5,$6,$7)", product.Title, product.Category, product.SubCategory, product.Description, product.Uid, product.InStock, product.Cost)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductDB) DeleteItem(Uid, ProductId string) error {
	_, err := pr.DB.Exec("UPDATE `Products` SET `IsDeleted` = 0 WHERE `id` = $1 AND `Uid` = $2", Uid, ProductId)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductDB) UpdateItem(product *Products, Pid int) error {
	_, err := pr.DB.Exec("UPDATE `Products` SET `title` = $1,`Category`= $2, `SubCategory`= $3, `Description`= $4, `Instock`= $5, `Cost` = $6 WHERE Uid=$7 AND id = $8", product.Title, product.Category, product.SubCategory, product.Description, product.InStock, product.Cost, product.Uid, Pid)
	if err != nil {
		return err
	}

	return nil
}
