package repository

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/hablof/product-registration/internal/database"
	"github.com/hablof/product-registration/internal/models"
	"github.com/hablof/product-registration/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	productsToAdd = []models.Product{
		{OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
		{OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
		{OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
	}
	productsToUpd = []models.Product{
		{OfferId: 1, Name: "name1", Price: 10, Quantity: 1},
		{OfferId: 2, Name: "name2", Price: 20, Quantity: 2},
		{OfferId: 3, Name: "name3", Price: 30, Quantity: 3},
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

	// это не тестирует методы репозитория...
	t.Run("добавляем записи для тестирования ProductsByFilter()", func(t *testing.T) {
		querystring, args, err := sq.Insert(tableName).
			Columns(sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol).
			Values(1, 1, "head", 10, 1).
			Values(1, 2, "body", 20, 0).
			Values(1, 3, "name1_3", 30, 3).
			Values(1, 4, "big changus", 40, 1).
			Values(1, 5, "Колесо", 1, 1).
			Values(1, 6, "Кросовок", 1, 1).
			Values(1, 7, "big melon", 2, 2).
			Values(1, 33, "head", 1, 1).
			Values(2, 1, "name2_1", 1, 1).
			Values(2, 3, "big boss", 3, 321).
			Values(2, 6, "Ветка", 1, 1).
			Values(2, 10, "big spoon", 1, 166).
			Values(3, 1, "name3_1", 1, 1).
			Values(3, 5, "Биткоинт", 1, 1).
			Values(3, 6, "big TV", 1, 1).
			Values(5, 5, "body", 2, 2).
			Values(6, 6, "submarine", 1, 1).
			Values(9, 9, "subwoofer", 2, 2).
			Values(15, 10, "name15_10", 71, 10).
			Values(20, 20, "subtitles", 1, 1).
			PlaceholderFormat(sq.Dollar).ToSql()

		if err != nil {
			assert.FailNow(t, err.Error())
		}

		rows, err := db.Exec(querystring, args...)
		if err != nil {
			assert.FailNow(t, err.Error())
		}

		rowsAffected, err := rows.RowsAffected()
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, int64(20), rowsAffected)
	})

	t.Run("проверка метода SellerProductIDs", func(t *testing.T) {

		testCases := map[uint64][]uint64{
			1:  {1, 2, 3, 4, 5, 6, 7, 33},
			2:  {1, 3, 6, 10},
			3:  {1, 5, 6},
			4:  {},
			5:  {5},
			6:  {6},
			7:  {},
			9:  {9},
			15: {10},
			20: {20},
		}

		for k, v := range testCases {
			offerIDs, err := r.SellerProductIDs(k)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			assert.Equal(t, v, offerIDs, fmt.Sprintf("SellerProductIDs with key \"%d\"", k))
		}
	})

	t.Run("выбираем по айди продавца", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{3, 9, 15},
			OfferIDs:  []uint64{},
			Substring: "",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 3, OfferId: 1, Name: "name3_1", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 5, Name: "Биткоинт", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 6, Name: "big TV", Price: 1, Quantity: 1},
				{SellerId: 9, OfferId: 9, Name: "subwoofer", Price: 2, Quantity: 2},
				{SellerId: 15, OfferId: 10, Name: "name15_10", Price: 71, Quantity: 10},
			},
			products)
	})

	t.Run("выбираем по айди продукта", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{},
			OfferIDs:  []uint64{1, 3},
			Substring: "",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 1, Name: "head", Price: 10, Quantity: 1},
				{SellerId: 1, OfferId: 3, Name: "name1_3", Price: 30, Quantity: 3},
				{SellerId: 2, OfferId: 1, Name: "name2_1", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 321},
				{SellerId: 3, OfferId: 1, Name: "name3_1", Price: 1, Quantity: 1},
			},
			products)
	})

	t.Run("выбираем по подстроке", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{},
			OfferIDs:  []uint64{},
			Substring: "big",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 40, Quantity: 1},
				{SellerId: 1, OfferId: 7, Name: "big melon", Price: 2, Quantity: 2},
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 321},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166},
				{SellerId: 3, OfferId: 6, Name: "big TV", Price: 1, Quantity: 1},
			},
			products)
	})

	t.Run("выбираем по айди продовца и айди продукта", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{1, 2},
			OfferIDs:  []uint64{3, 4},
			Substring: "",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 3, Name: "name1_3", Price: 30, Quantity: 3},
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 40, Quantity: 1},
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 321},
			},
			products)
	})

	t.Run("выбираем по айди продовца и по подстроке", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{2, 3, 4},
			OfferIDs:  []uint64{},
			Substring: "big",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 321},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166},
				{SellerId: 3, OfferId: 6, Name: "big TV", Price: 1, Quantity: 1},
			},
			products)
	})

	t.Run("выбираем по айди продукта и по подстроке", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{},
			OfferIDs:  []uint64{4, 10},
			Substring: "big",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 40, Quantity: 1},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166},
			},
			products)
	})

	t.Run("выбираем ипо айди продовца, и по айди продукта, и по подстроке", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{},
			OfferIDs:  []uint64{4, 10},
			Substring: "big",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 40, Quantity: 1},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166},
			},
			products)
	})

	t.Run("выбираем без фильтра", func(t *testing.T) {
		products, err := r.ProductsByFilter(service.RequestFilter{
			SellerIDs: []uint64{},
			OfferIDs:  []uint64{},
			Substring: "",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t,
			[]models.Product{
				{SellerId: 1, OfferId: 1, Name: "head", Price: 10, Quantity: 1},
				{SellerId: 1, OfferId: 2, Name: "body", Price: 20, Quantity: 0},
				{SellerId: 1, OfferId: 3, Name: "name1_3", Price: 30, Quantity: 3},
				{SellerId: 1, OfferId: 4, Name: "big changus", Price: 40, Quantity: 1},
				{SellerId: 1, OfferId: 5, Name: "Колесо", Price: 1, Quantity: 1},
				{SellerId: 1, OfferId: 6, Name: "Кросовок", Price: 1, Quantity: 1},
				{SellerId: 1, OfferId: 7, Name: "big melon", Price: 2, Quantity: 2},
				{SellerId: 1, OfferId: 33, Name: "head", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 1, Name: "name2_1", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 3, Name: "big boss", Price: 3, Quantity: 321},
				{SellerId: 2, OfferId: 6, Name: "Ветка", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 10, Name: "big spoon", Price: 1, Quantity: 166},
				{SellerId: 3, OfferId: 1, Name: "name3_1", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 5, Name: "Биткоинт", Price: 1, Quantity: 1},
				{SellerId: 3, OfferId: 6, Name: "big TV", Price: 1, Quantity: 1},
				{SellerId: 5, OfferId: 5, Name: "body", Price: 2, Quantity: 2},
				{SellerId: 6, OfferId: 6, Name: "submarine", Price: 1, Quantity: 1},
				{SellerId: 9, OfferId: 9, Name: "subwoofer", Price: 2, Quantity: 2},
				{SellerId: 15, OfferId: 10, Name: "name15_10", Price: 71, Quantity: 10},
				{SellerId: 20, OfferId: 20, Name: "subtitles", Price: 1, Quantity: 1},
			},
			products)
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
