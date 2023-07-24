package repository

import (
	"errors"
	"log"
	"testing"

	"github.com/hablof/product-registration/internal/models"
	"github.com/hablof/product-registration/internal/service"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRepository_ManageProducts(t *testing.T) {
	db, mockCtrl, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name             string
		sellerId         uint64
		productsToAdd    []models.Product
		productsToDelete []models.Product
		productsToUpdate []models.Product
		mockBehaviour    func(m sqlxmock.Sqlmock)
		wantErr          error
	}{
		{
			name:             "empty request",
			sellerId:         42,
			productsToAdd:    []models.Product{},
			productsToDelete: []models.Product{},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
			},
			wantErr: ErrEmptyRequest,
		},
		{
			name:             "transaction start failed",
			sellerId:         42,
			productsToAdd:    []models.Product{},
			productsToDelete: []models.Product{{OfferId: 1, Name: "name1", Price: 1, Quantity: 1}},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(errors.New("cannot start tx"))
			},
			wantErr: ErrTxFailed,
		},
		{
			name:             "insert query exec failed",
			sellerId:         42,
			productsToAdd:    []models.Product{{OfferId: 1, Name: "name1", Price: 1, Quantity: 1}},
			productsToDelete: []models.Product{},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO products (seller_id,offer_id,name,price,quantity) 
				VALUES ($1,$2,$3,$4,$5) ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
				name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity`).
					WillReturnError(errors.New("exec error"))
			},
			wantErr: ErrQueryExecFailed,
		},
		{
			name:             "insert query exec result with error",
			sellerId:         42,
			productsToAdd:    []models.Product{{OfferId: 1, Name: "name1", Price: 1, Quantity: 1}},
			productsToDelete: []models.Product{},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO products (seller_id,offer_id,name,price,quantity) 
				VALUES ($1,$2,$3,$4,$5) ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
				name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity`).
					WillReturnResult(sqlxmock.NewErrorResult(errors.New("result exec err")))
			},
			wantErr: ErrQueryExecFailed,
		},
		{
			name:             "delete query exec failed",
			sellerId:         42,
			productsToAdd:    []models.Product{},
			productsToDelete: []models.Product{{OfferId: 1, Name: "name1", Price: 1, Quantity: 1}},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("DELETE FROM products WHERE offer_id IN ($1) AND seller_id = $2").
					WillReturnError(errors.New("exec error"))
			},
			wantErr: ErrQueryExecFailed,
		},
		{
			name:             "delete query exec result with error",
			sellerId:         42,
			productsToAdd:    []models.Product{},
			productsToDelete: []models.Product{{OfferId: 1, Name: "name1", Price: 1, Quantity: 1}},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("DELETE FROM products WHERE offer_id IN ($1) AND seller_id = $2").
					WillReturnResult(sqlxmock.NewErrorResult(errors.New("result exec err")))
			},
			wantErr: ErrQueryExecFailed,
		},
		{
			name:          "request to delete only",
			sellerId:      42,
			productsToAdd: []models.Product{},
			productsToDelete: []models.Product{
				{OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("DELETE FROM products WHERE offer_id IN ($1,$2,$3) AND seller_id = $4").
					WithArgs(1, 2, 3, 42).
					WillReturnResult(sqlxmock.NewResult(0, 3))
				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name:          "request to update only",
			sellerId:      42,
			productsToAdd: []models.Product{},
			productsToUpdate: []models.Product{
				{OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			productsToDelete: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO products (seller_id,offer_id,name,price,quantity) 
					VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10),($11,$12,$13,$14,$15) 
					ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
					name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity`).
					WithArgs(42, 1, "name1", 1, 1, 42, 2, "name2", 2, 2, 42, 3, "name3", 3, 3).
					WillReturnResult(sqlxmock.NewResult(0, 3))
				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name:     "request to insert only",
			sellerId: 42,
			productsToAdd: []models.Product{
				{OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			productsToDelete: []models.Product{},
			productsToUpdate: []models.Product{},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO products (seller_id,offer_id,name,price,quantity) 
					VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10),($11,$12,$13,$14,$15) 
					ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
					name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity`).
					WithArgs(42, 1, "name1", 1, 1, 42, 2, "name2", 2, 2, 42, 3, "name3", 3, 3).
					WillReturnResult(sqlxmock.NewResult(0, 3))
				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name:     "request to 1 insert, 1 update, 1 delete",
			sellerId: 42,
			productsToAdd: []models.Product{
				{OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
			},
			productsToDelete: []models.Product{
				{OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
			},
			productsToUpdate: []models.Product{
				{OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(`INSERT INTO products (seller_id,offer_id,name,price,quantity) 
					VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10) 
					ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
					name = EXCLUDED.name, price = EXCLUDED.price, quantity = EXCLUDED.quantity`).
					WithArgs(42, 1, "name1", 1, 1, 42, 3, "name3", 3, 3).
					WillReturnResult(sqlxmock.NewResult(0, 2))
				m.ExpectExec("DELETE FROM products WHERE offer_id IN ($1) AND seller_id = $2").
					WithArgs(2, 42).
					WillReturnResult(sqlxmock.NewResult(0, 1))
				m.ExpectCommit()
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db)
			tt.mockBehaviour(mockCtrl)

			err := r.ManageProducts(tt.sellerId, tt.productsToAdd, tt.productsToDelete, tt.productsToUpdate)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRepository_ProductsByFilter(t *testing.T) {
	db, mockCtrl, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		filter        service.RequestFilter
		mockBehaviour func(m sqlxmock.Sqlmock)
		want          []models.Product
		wantErr       error
	}{
		{
			name: "error query execution",
			filter: service.RequestFilter{
				SellerIDs: []uint64{1, 2, 3},
				OfferIDs:  []uint64{1},
				Substring: "err",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT`
				m.ExpectQuery(reg).WithArgs(1, 2, 3, 1, `%err%`).WillReturnError(errors.New("error query execution"))
			},
			want:    nil,
			wantErr: ErrQueryExecFailed,
		},
		{
			name: "filter by seller id",
			filter: service.RequestFilter{
				SellerIDs: []uint64{1, 2, 3},
				OfferIDs:  []uint64{},
				Substring: "",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT.+seller_id IN`
				rows := sqlxmock.NewRows([]string{sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol}).
					AddRow(1, 1, "head", 1, 1).
					AddRow(1, 2, "body", 2, 2).
					AddRow(1, 3, "name1_3", 3, 3).
					AddRow(2, 1, "name2_1", 1, 1).
					AddRow(3, 1, "name3_1", 1, 1)
				m.ExpectQuery(reg).WithArgs(1, 2, 3).WillReturnRows(rows)
			},
			want: []models.Product{
				{SellerId: 1, OfferId: 1, Name: "head", Price: 1, Quantity: 1},
				{SellerId: 1, OfferId: 2, Name: "body", Price: 2, Quantity: 2},
				{SellerId: 1, OfferId: 3, Name: "name1_3", Price: 3, Quantity: 3},
				{SellerId: 2, OfferId: 1, Name: "name2_1", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 1, Name: "name3_1", Price: 1, Quantity: 1},
			},
			wantErr: nil,
		},
		{
			name: "filter by offer id",
			filter: service.RequestFilter{
				SellerIDs: []uint64{},
				OfferIDs:  []uint64{1, 5, 10},
				Substring: "",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT.+offer_id IN`
				rows := sqlxmock.NewRows([]string{sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol}).
					AddRow(1, 1, "head", 1, 1).
					AddRow(5, 5, "body", 2, 2).
					AddRow(15, 10, "name15_10", 1, 1)
				m.ExpectQuery(reg).WithArgs(1, 5, 10).WillReturnRows(rows)
			},
			want: []models.Product{
				{SellerId: 1, OfferId: 1, Name: "head", Price: 1, Quantity: 1},
				{SellerId: 5, OfferId: 5, Name: "body", Price: 2, Quantity: 2},
				{SellerId: 15, OfferId: 10, Name: "name15_10", Price: 1, Quantity: 1},
			},
			wantErr: nil,
		},
		{
			name: "filter by substring",
			filter: service.RequestFilter{
				SellerIDs: []uint64{},
				OfferIDs:  []uint64{},
				Substring: "sub",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT.+name LIKE`
				rows := sqlxmock.NewRows([]string{sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol}).
					AddRow(6, 6, "submarine", 1, 1).
					AddRow(9, 9, "subwoofer", 2, 2).
					AddRow(20, 20, "subtitles", 1, 1)
				m.ExpectQuery(reg).WithArgs(`%sub%`).WillReturnRows(rows)
			},
			want: []models.Product{
				{SellerId: 6, OfferId: 6, Name: "submarine", Price: 1, Quantity: 1},
				{SellerId: 9, OfferId: 9, Name: "subwoofer", Price: 2, Quantity: 2},
				{SellerId: 20, OfferId: 20, Name: "subtitles", Price: 1, Quantity: 1},
			},
			wantErr: nil,
		},
		{
			name: "filter by seller_id, offer_id",
			filter: service.RequestFilter{
				SellerIDs: []uint64{1, 2, 3},
				OfferIDs:  []uint64{5, 6},
				Substring: "",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT.+seller_id IN.+offer_id IN`
				rows := sqlxmock.NewRows([]string{sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol}).
					AddRow(1, 5, "Колесо", 1, 1).
					AddRow(1, 6, "Кросовок", 1, 1).
					AddRow(2, 6, "Ветка", 1, 1).
					AddRow(3, 5, "Биткоинт", 1, 1)
				m.ExpectQuery(reg).WithArgs(1, 2, 3, 5, 6).WillReturnRows(rows)
			},
			want: []models.Product{
				{SellerId: 1, OfferId: 5, Name: "Колесо", Price: 1, Quantity: 1},
				{SellerId: 1, OfferId: 6, Name: "Кросовок", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 6, Name: "Ветка", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 5, Name: "Биткоинт", Price: 1, Quantity: 1},
			},
			wantErr: nil,
		},
		{
			name: "filter by seller_id, offer_id, name",
			filter: service.RequestFilter{
				SellerIDs: []uint64{1, 2, 3},
				OfferIDs:  []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				Substring: "big",
			},
			mockBehaviour: func(m sqlxmock.Sqlmock) {
				reg := `SELECT.+(?:(?:seller_id IN|offer_id IN|name LIKE).+){3}` // ровно три раза встретим одно из ...
				rows := sqlxmock.NewRows([]string{sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol}).
					AddRow(1, 4, "big changus", 1, 1).
					AddRow(1, 7, "big melon", 2, 2).
					AddRow(2, 3, "big boss", 3, 3).
					AddRow(2, 10, "big spoon", 1, 1).
					AddRow(3, 6, "big TV", 1, 1)
				m.ExpectQuery(reg).WithArgs(1, 2, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, `%big%`).WillReturnRows(rows)
			},
			want: []models.Product{
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 1, Quantity: 1},
				{SellerId: 1, OfferId: 7, Name: "big melon", Price: 2, Quantity: 2},
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 3},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 6, Name: "big TV", Price: 1, Quantity: 1},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db)
			tt.mockBehaviour(mockCtrl)

			log.Println(tt.name)

			products, err := r.ProductsByFilter(tt.filter)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, products)
		})
	}
}
