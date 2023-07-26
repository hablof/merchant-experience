package xlsxparser

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hablof/product-registration/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestXLSXparser(t *testing.T) {
	testCases := []struct {
		testname        string
		fileName        string
		want            []models.ProductUpdate
		wantProductErrs []error
		wantErr         error
	}{
		{
			testname:        "empty file",
			fileName:        "example_empty.xlsx",
			want:            nil,
			wantProductErrs: nil,
			wantErr:         ErrEmptySheet,
		},
		{
			testname:        "duplicates",
			fileName:        "example_duplicates.xlsx",
			want:            nil,
			wantProductErrs: nil,
			wantErr:         ErrHasDuplicates,
		},
		{
			testname:        "",
			fileName:        "example_with_invalid_offer_id_col.xlsx",
			want:            nil,
			wantProductErrs: nil,
			wantErr:         ErrInvalidIDs,
		},
		{
			testname: "file with errors",
			fileName: "example_with_errors.xlsx",
			want: []models.ProductUpdate{
				{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true},
				{Product: models.Product{OfferId: 9, Name: "name2_1", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 10, Name: "big boss", Price: 3, Quantity: 321}, Available: true},
				{Product: models.Product{OfferId: 11, Name: "Ветка", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 12, Name: "big spoon", Price: 1, Quantity: 166}, Available: true},
				{Product: models.Product{OfferId: 13, Name: "name3_1", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 14, Name: "Биткоинт", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 15, Name: "big TV", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 16, Name: "body", Price: 2, Quantity: 2}, Available: true},
				{Product: models.Product{OfferId: 17, Name: "submarine", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 18, Name: "subwoofer", Price: 2, Quantity: 2}, Available: true},
				{Product: models.Product{OfferId: 19, Name: "name15_10", Price: 71, Quantity: 10}, Available: true},
				{Product: models.Product{OfferId: 20, Name: "subtitles", Price: 1, Quantity: 1}, Available: true},
			},
			wantProductErrs: []error{
				ErrProductParsing{
					Row:    3,
					Field:  "name",
					ErrMsg: models.MsgTooLongName,
				},
				ErrProductParsing{
					Row:   4,
					Field: "price",
					ErrMsg: (&strconv.NumError{
						Func: "ParseUint",
						Num:  "0-40",
						Err:  strconv.ErrSyntax,
					}).Error(),
				},
				ErrProductParsing{
					Row:   5,
					Field: "price",
					ErrMsg: (&strconv.NumError{
						Func: "ParseUint",
						Num:  "-666",
						Err:  strconv.ErrSyntax,
					}).Error(),
				},
				ErrProductParsing{
					Row:   6,
					Field: "quantity",
					ErrMsg: (&strconv.NumError{
						Func: "ParseUint",
						Num:  "0-40",
						Err:  strconv.ErrSyntax,
					}).Error(),
				},
				ErrProductParsing{
					Row:   7,
					Field: "quantity",
					ErrMsg: (&strconv.NumError{
						Func: "ParseUint",
						Num:  "-666",
						Err:  strconv.ErrSyntax,
					}).Error(),
				},
				ErrProductParsing{
					Row:   8,
					Field: "available",
					ErrMsg: (&strconv.NumError{
						Func: "ParseBool",
						Num:  "абра-кадабра",
						Err:  strconv.ErrSyntax,
					}).Error(),
				},
			},
			wantErr: nil,
		},
		{
			testname: "20 rows",
			fileName: "example1.xlsx",
			want: []models.ProductUpdate{
				{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true},
				{Product: models.Product{OfferId: 3, Name: "name1_3", Price: 30, Quantity: 3}, Available: true},
				{Product: models.Product{OfferId: 4, Name: "big changus", Price: 40, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 5, Name: "Колесо", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 6, Name: "Кросовок", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 7, Name: "big melon", Price: 2, Quantity: 2}, Available: false},
				{Product: models.Product{OfferId: 8, Name: "head", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 9, Name: "name2_1", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 10, Name: "big boss", Price: 3, Quantity: 321}, Available: false},
				{Product: models.Product{OfferId: 11, Name: "Ветка", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 12, Name: "big spoon", Price: 1, Quantity: 166}, Available: true},
				{Product: models.Product{OfferId: 13, Name: "name3_1", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 14, Name: "Биткоинт", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 15, Name: "big TV", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 16, Name: "body", Price: 2, Quantity: 2}, Available: false},
				{Product: models.Product{OfferId: 17, Name: "submarine", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 18, Name: "subwoofer", Price: 2, Quantity: 2}, Available: true},
				{Product: models.Product{OfferId: 19, Name: "name15_10", Price: 71, Quantity: 10}, Available: true},
				{Product: models.Product{OfferId: 20, Name: "subtitles", Price: 1, Quantity: 1}, Available: true},
			},
			wantProductErrs: nil,
			wantErr:         nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.testname, func(t *testing.T) {
			filename := filepath.Join("test", tt.fileName)
			f, err := os.Open(filename)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			p := NewParser()

			parsedProducts, parseErrs, err := p.ParseProducts(f)
			assert.Equal(t, tt.wantErr, err, "method errors")
			assert.Equal(t, tt.wantProductErrs, parseErrs, "parse errors")
			assert.Equal(t, tt.want, parsedProducts, "parsed products")
		})
	}
}
