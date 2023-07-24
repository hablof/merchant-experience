package xslxparser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hablof/product-registration/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestXLSXparser(t *testing.T) {
	testCases := []struct {
		testname string
		fileName string
		want     []models.ProductUpdate
		wantErr  error
	}{
		{
			testname: "empty file",
			fileName: "example_empty.xlsx",
			want:     nil,
			wantErr:  ErrEmptySheet,
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
				{Product: models.Product{OfferId: 33, Name: "head", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 1, Name: "name2_1", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 3, Name: "big boss", Price: 3, Quantity: 321}, Available: false},
				{Product: models.Product{OfferId: 6, Name: "Ветка", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166}, Available: true},
				{Product: models.Product{OfferId: 1, Name: "name3_1", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 5, Name: "Биткоинт", Price: 1, Quantity: 1}, Available: false},
				{Product: models.Product{OfferId: 6, Name: "big TV", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 5, Name: "body", Price: 2, Quantity: 2}, Available: false},
				{Product: models.Product{OfferId: 6, Name: "submarine", Price: 1, Quantity: 1}, Available: true},
				{Product: models.Product{OfferId: 9, Name: "subwoofer", Price: 2, Quantity: 2}, Available: true},
				{Product: models.Product{OfferId: 10, Name: "name15_10", Price: 71, Quantity: 10}, Available: true},
				{Product: models.Product{OfferId: 20, Name: "subtitles", Price: 1, Quantity: 1}, Available: true},
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.testname, func(t *testing.T) {
			filename := filepath.Join("test", tt.fileName)
			f, err := os.Open(filename)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			parsedProducts, err := ParseProducts(f)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, parsedProducts)
		})
	}
}
