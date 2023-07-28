package xlsxparser

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/hablof/merchant-experience/internal/models"
	"github.com/xuri/excelize/v2"
)

var (
	ErrEmptyDoc      = errors.New("empty document")
	ErrEmptySheet    = errors.New("empty sheet")
	ErrHasDuplicates = errors.New("sheet contain offer_id duplicates")
	ErrInvalidIDs    = errors.New("offer_id column has invalid value(s)")
	ErrFailedToRead  = errors.New("cannot read document")
)

type ErrProductParsing struct {
	Row    uint64 `json:"row"`
	Field  string `json:"field"`
	ErrMsg string `json:"errMsg"`
}

func (e ErrProductParsing) Error() string {
	return fmt.Sprintf("product invalid: id=%d, field=%s, err=%s", e.Row, e.Field, e.ErrMsg)
}

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

// метод не знает ничего про seller_id
func (Parser) ParseProducts(r io.Reader) (productUpdates []models.ProductUpdate, productErrs []error, methodErr error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		log.Println(err)
		return nil, nil, ErrFailedToRead
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	rows, err := prepare(f)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// основной цикл
	productUpdates = make([]models.ProductUpdate, 0, len(rows))
	productErrs = make([]error, 0, len(rows))
	for rowNumber, row := range rows {

		productUnit := models.Product{}
		updateUnit := models.ProductUpdate{}
		isValid := true
		// cols:
		// [0] offer_id  - уникальный идентификатор товара в системе продавца
		// [1] name      - название товара
		// [2] price     - цена в рублях
		// [3] quantity  - количество товара на складе продавца
		// [4] available - true/false, в случае false продавец хочет удалить товар из нашей базы

		// парсим offer_id
		offerId, err := strconv.ParseUint(row[0], 10, 64)
		if err != nil {
			isValid = false
			e := ErrProductParsing{
				Row:    uint64(rowNumber + 1), // человеческий счёт
				Field:  "offer_id",
				ErrMsg: err.Error(),
			}
			productErrs = append(productErrs, e)
		}

		// обрезаем пробелы у name
		name := strings.TrimSpace(row[1])

		// парсим price
		price, err := strconv.ParseUint(row[2], 10, 64)
		if err != nil {
			isValid = false
			e := ErrProductParsing{
				Row:    uint64(rowNumber + 1), // человеческий счёт
				Field:  "price",
				ErrMsg: err.Error(),
			}
			productErrs = append(productErrs, e)
		}

		// парсим quantity
		quantity, err := strconv.ParseUint(row[3], 10, 64)
		if err != nil {
			isValid = false
			e := ErrProductParsing{
				Row:    uint64(rowNumber + 1), // человеческий счёт
				Field:  "quantity",
				ErrMsg: err.Error(),
			}
			productErrs = append(productErrs, e)
		}

		// парсим available
		available, err := strconv.ParseBool(row[4])
		if err != nil {
			isValid = false
			e := ErrProductParsing{
				Row:    uint64(rowNumber + 1), // человеческий счёт
				Field:  "available",
				ErrMsg: err.Error(),
			}
			productErrs = append(productErrs, e)
		}

		productUnit.OfferId = offerId
		productUnit.Name = name
		productUnit.Price = price
		productUnit.Quantity = quantity

		// валидируем по логике домена
		var validationErr models.ErrProductValidation
		if err := productUnit.Validate(); errors.As(err, &validationErr) {
			isValid = false
			e := ErrProductParsing{
				Row:    uint64(rowNumber + 1), // человеческий счёт
				Field:  validationErr.Field,
				ErrMsg: validationErr.ErrMsg,
			}
			productErrs = append(productErrs, e)
		}

		updateUnit.Product = productUnit
		updateUnit.Available = available

		if isValid {
			productUpdates = append(productUpdates, updateUnit)
		}
	}

	productUpdates = append(make([]models.ProductUpdate, 0, len(productUpdates)), productUpdates...)
	productErrs = append(make([]error, 0, len(productErrs)), productErrs...)
	if len(productErrs) == 0 {
		productErrs = nil
	}

	return productUpdates, productErrs, nil
}

func prepare(f *excelize.File) ([][]string, error) {
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

	cols, err := f.GetCols(sheetList[0])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	offerIDs := make([]uint64, 0, len(cols[0]))
	for _, str := range cols[0] {
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			log.Println(err)
			return nil, ErrInvalidIDs // человеческая система счёта
		}

		offerIDs = append(offerIDs, u)
	}

	// check for duplicates
	if hasDuplicates(offerIDs) {
		log.Println("sheet has offerID duplicates")
		return nil, ErrHasDuplicates
	}
	return rows, nil
}

func hasDuplicates[T comparable](slice []T) bool {
	m := make(map[T]struct{}, len(slice))

	for _, elem := range slice {
		if _, ok := m[elem]; ok {
			return true
		} else {
			m[elem] = struct{}{}
		}
	}

	return false
}
