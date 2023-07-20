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
	OfferId  uint64
	Name     string
	Price    uint64
	Quantity uint64
}

func (p Product) Validate() error {
	if utf8.RuneCountInString(p.Name) > 100 {
		return ErrTooLongName
	}

	return nil
}
