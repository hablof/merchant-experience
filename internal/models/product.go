package models

type ProductUpdate struct {
	Product Product
	Available bool
}

type Product struct {
	OfferId uint64
	Name string
	Price uint64
	Quantity uint64
}