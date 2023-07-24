package xslxparser

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/hablof/product-registration/internal/models"
	"github.com/xuri/excelize/v2"
)

var (
	ErrEmptyDoc   = errors.New("empty document")
	ErrEmptySheet = errors.New("empty sheet")
)

// метод не знает ничего про seller_id
func ParseProducts(r io.Reader) ([]models.ProductUpdate, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		log.Println("empty document")
		return nil, ErrEmptyDoc
	}

	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(rows) == 0 {
		log.Println("empty sheet")
		return nil, ErrEmptySheet
	}

	productUpdates := make([]models.ProductUpdate, 0, len(rows))
	for _, row := range rows {
		productUnit := models.Product{}
		updateUnit := models.ProductUpdate{}

		// cols:
		// [0] offer_id  - уникальный идентификатор товара в системе продавца
		// [1] name      - название товара
		// [2] price     - цена в рублях
		// [3] quantity  - количество товара на складе продавца
		// [4] available - true/false, в случае false продавец хочет удалить товар из нашей базы

		offerId, err := strconv.ParseUint(row[0], 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		name := strings.TrimSpace(row[1])

		price, err := strconv.ParseUint(row[2], 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		quantity, err := strconv.ParseUint(row[3], 10, 64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		available, err := strconv.ParseBool(row[4])
		if err != nil {
			log.Println(err)
			return nil, err
		}

		productUnit.OfferId = offerId
		productUnit.Name = name
		productUnit.Price = price
		productUnit.Quantity = quantity

		updateUnit.Product = productUnit
		updateUnit.Available = available

		productUpdates = append(productUpdates, updateUnit)
	}

	return productUpdates, nil
}
