package models

import (
	"fmt"
	"unicode/utf8"
)

const (
	MsgTooLongName = "too long name"
)

type ErrProductValidation struct {
	OfferId uint64 `json:"offerId"`
	Field   string `json:"field"`
	ErrMsg  string `json:"errMsg"`
}

func (e ErrProductValidation) Error() string {
	return fmt.Sprintf("product invalid: id=%d, field=%s, err=%s", e.OfferId, e.Field, e.ErrMsg)
}

type ProductUpdate struct {
	Product Product
	// SellerId  uint64
	Available bool
}

type Product struct {
	SellerId uint64 `db:"seller_id" json:"sellerId"`
	OfferId  uint64 `db:"offer_id"  json:"offerId"`
	Name     string `db:"name"      json:"name"`
	Price    uint64 `db:"price"     json:"price"`
	Quantity uint64 `db:"quantity"  json:"quantity"`
}

// returns ErrProductValidation type
func (p Product) Validate() error {
	e := ErrProductValidation{
		OfferId: p.OfferId,
	}

	switch {
	case utf8.RuneCountInString(p.Name) > 100:
		e.Field = "name"
		e.ErrMsg = MsgTooLongName
		return e
	}

	return nil
}
