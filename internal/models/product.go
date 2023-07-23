package models

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrTooLongName = errors.New("too long name")
)

type ProductUpdate struct {
	Product Product
	// SellerId  uint64
	Available bool
}

type Product struct {
	SellerId uint64 `db:"seller_id"`
	OfferId  uint64 `db:"offer_id"`
	Name     string `db:"name"`
	Price    uint64 `db:"price"`
	Quantity uint64 `db:"quantity"`
}

func (p Product) Validate() error {
	if utf8.RuneCountInString(p.Name) > 100 {
		return ErrTooLongName
	}

	return nil
}
