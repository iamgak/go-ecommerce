package main

import (
	"ecommerce.iamgak.com/internals"
	"fmt"
)

type Application struct {
	IsAuthenticated bool
	UserType        string
	Product         *products.Product
}

func main() {
	fmt.Print("hello")

}
