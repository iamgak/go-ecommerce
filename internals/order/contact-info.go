package order

type UserInterf interface {
	CreateContactInfo(*ContactInfo) error
	UpdateContactInfo(*ContactInfo, int) error
}

type ContactInfo struct {
	OrderId int
	Addr    string
	Phone   string
	Uid     int
}
