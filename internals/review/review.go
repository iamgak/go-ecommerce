package review

import "database/sql"

type ReviewFunc interface {
	Review(*Review) error
}

type Review struct {
	CartId     int
	Addr       string
	Quantity   int
	TotalPrice float32
}

type ReviewDB struct {
	DB *sql.DB
}

func (r *Review) Review(review *Review) error {

	return nil
}
