package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("errors: invalid credentials")
	ErrDuplicateEmail     = errors.New("errors: duplicate email")
	ErrNoRecord           = errors.New("errors: no matching record found")
	ErrCantFindProduct    = errors.New("errors: can't find product")
	ErrCantDecodeProducts = errors.New("errors: can't find product")
	ErrUserIDIsNotValid   = errors.New("errors: user is not valid")
	ErrCantAddInCart      = errors.New("errors: cannot add product to cart")
	ErrCantRemoveItem     = errors.New("errors: cannot remove item from cart")
	ErrCantGetItem        = errors.New("errors: cannot get item from cart ")
	ErrCantBuyCartItem    = errors.New("errors: cannot update the purchase")
)
