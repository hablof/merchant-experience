package repository

import (
	"testing"

	"github.com/hablof/product-registration/internal/database"
	"github.com/hablof/product-registration/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	productsToAdd = []models.Product{
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
	}
	productsToUpd = []models.Product{
		{
			OfferId:  1,
			Name:     "name1",
			Price:    10,
			Quantity: 1,
		},
		{
			OfferId:  2,
			Name:     "name2",
			Price:    20,
			Quantity: 2,
		},
		{
			OfferId:  3,
			Name:     "name3",
			Price:    30,
			Quantity: 3,
		},
	}
)

func TestRepository(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	db, err := database.NewPostgres()
	if !assert.NoError(t, err) {
		assert.FailNow(t, "no database connection")
	}

	r := NewRepository(db)

	defer teardown(t, db)
	setup(t, db)

	t.Run("добавляем три записи", func(t *testing.T) {
		if err := r.ManageProducts(0, productsToAdd, nil, nil); err != nil {
			assert.FailNow(t, err.Error())
		}

		result, err := db.Queryx(`SELECT offer_id, name, price, quantity FROM products;`)
		if err != nil {
			assert.FailNow(t, err.Error())
		}

		results := make([]models.Product, 0, 3)
		for result.Next() {
			unit := models.Product{}
			if err := result.StructScan(&unit); err != nil {
				assert.Fail(t, err.Error())
			}
			results = append(results, unit)
		}
		assert.Equal(t, productsToAdd, results)
	})

	t.Run("меняем все три добавленные записи", func(t *testing.T) {
		if err := r.ManageProducts(0, nil, nil, productsToUpd); err != nil {
			assert.FailNow(t, err.Error())
		}

		result, err := db.Queryx(`SELECT offer_id, name, price, quantity FROM products;`)
		if err != nil {
			assert.FailNow(t, err.Error())
		}

		results := make([]models.Product, 0, 3)
		for result.Next() {
			unit := models.Product{}
			if err := result.StructScan(&unit); err != nil {
				assert.Fail(t, err.Error())
			}
			results = append(results, unit)
		}
		assert.Equal(t, productsToUpd, results)
	})

	t.Run("удаляем все три записи", func(t *testing.T) {
		productsToDel := productsToUpd
		if err := r.ManageProducts(0, nil, productsToDel, nil); err != nil {
			assert.FailNow(t, err.Error())
		}

		count := 0
		if err := db.QueryRowx(`SELECT COUNT(*) FROM products;`).Scan(&count); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.Equal(t, 0, count)
	})

}

func teardown(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS products;`)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}

func setup(t *testing.T, db *sqlx.DB) {

	if _, err := db.Exec(`DROP TABLE IF EXISTS products;`); err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err := db.Exec(`CREATE TABLE products (
		seller_id BIGINT      NOT NULL,
		offer_id  BIGINT      NOT NULL, 
		name      VARCHAR(100) NOT NULL,
		price     BIGINT      NOT NULL,
		quantity  BIGINT      NOT NULL,
		PRIMARY KEY(seller_id, offer_id),
		CONSTRAINT no_duplicates UNIQUE(seller_id, offer_id)
	);`)

	if err != nil {
		assert.FailNow(t, err.Error())
	}
}
