package repository

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/hablof/product-registration/internal/models"
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
			name:          "request to delete only",
			sellerId:      42,
			productsToAdd: []models.Product{},
			productsToDelete: []models.Product{
				{
					OfferId:  1,
					Name:     "name1",
					Price:    1,
					Quantity: 1,
				},
				{
					OfferId:  2,
					Name:     "name2",
					Price:    2,
					Quantity: 2,
				},
				{
					OfferId:  3,
					Name:     "name3",
					Price:    3,
					Quantity: 3,
				},
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
				{
					OfferId:  1,
					Name:     "name1",
					Price:    1,
					Quantity: 1,
				},
				{
					OfferId:  2,
					Name:     "name2",
					Price:    2,
					Quantity: 2,
				},
				{
					OfferId:  3,
					Name:     "name3",
					Price:    3,
					Quantity: 3,
				},
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
				{
					OfferId:  1,
					Name:     "name1",
					Price:    1,
					Quantity: 1,
				},
				{
					OfferId:  2,
					Name:     "name2",
					Price:    2,
					Quantity: 2,
				},
				{
					OfferId:  3,
					Name:     "name3",
					Price:    3,
					Quantity: 3,
				},
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
				{
					OfferId:  1,
					Name:     "name1",
					Price:    1,
					Quantity: 1,
				},
			},
			productsToDelete: []models.Product{
				{
					OfferId:  2,
					Name:     "name2",
					Price:    2,
					Quantity: 2,
				},
			},
			productsToUpdate: []models.Product{
				{
					OfferId:  3,
					Name:     "name3",
					Price:    3,
					Quantity: 3,
				},
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
			r := Repository{
				db:        db,
				initQuery: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
			}
			tt.mockBehaviour(mockCtrl)
			err := r.ManageProducts(tt.sellerId, tt.productsToAdd, tt.productsToDelete, tt.productsToUpdate)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}
