package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hablof/product-registration/internal/database"
	"github.com/hablof/product-registration/internal/gateway"
	"github.com/hablof/product-registration/internal/repository"
	"github.com/hablof/product-registration/internal/router"
	"github.com/hablof/product-registration/internal/service"
	"github.com/hablof/product-registration/internal/xlsxparser"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const (
	serverHostPort     = ":8015"
	tableReqHostPort   = "http://127.0.0.1:8015"
	respBodyWithErrors = `{"added":14,"updated":0,"deleted":0,"errors":[{"row":3,"field":"name","errMsg":"too long name"},{"row":4,"field":"price","errMsg":"strconv.ParseUint: parsing \"0-40\": invalid syntax"},{"row":5,"field":"price","errMsg":"strconv.ParseUint: parsing \"-666\": invalid syntax"},{"row":6,"field":"quantity","errMsg":"strconv.ParseUint: parsing \"0-40\": invalid syntax"},{"row":7,"field":"quantity","errMsg":"strconv.ParseUint: parsing \"-666\": invalid syntax"},{"row":8,"field":"available","errMsg":"strconv.ParseBool: parsing \"абра-кадабра\": invalid syntax"}]}`
)

type postTableJson struct {
	TableURL string `json:"tableURL"`
	SellerId uint64 `json:"sellerId"`
}

func databaseTeardown(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS products;`)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}

func databaseSetup(t *testing.T, db *sqlx.DB) {

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

func TestMicroservice(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// setup
	fileServer := &http.Server{
		Addr:        serverHostPort,
		Handler:     newTestServerHandler(),
		ReadTimeout: 1 * time.Second,
	}
	go func(testServer *http.Server) {
		if err := testServer.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			assert.FailNow(t, err.Error())
		}
	}(fileServer)
	defer fileServer.Close()

	db, err := database.NewPostgres()
	if !assert.NoError(t, err) {
		assert.FailNow(t, "no database connection")
	}

	r := repository.NewRepository(db)
	s := service.NewService(r)
	g := gateway.NewGateway()
	p := xlsxparser.NewParser()
	handler := router.NewRouter(s, g, p)

	databaseSetup(t, db)
	defer databaseTeardown(t, db)

	name := "post table with duplicates"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/xlsxparser/test/example_duplicates.xlsx", 42)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, "xslx file has duplicates", w.Body.String())
	})

	name = "post empty table"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/xlsxparser/test/example_empty.xlsx", 42)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, "bad xslx file", w.Body.String())
	})

	name = "post table with invalid offer_id column"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/xlsxparser/test/example_with_invalid_offer_id_col.xlsx", 42)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, "offer_id column has invalid value(s)", w.Body.String())
	})

	name = "post txt file"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/testtables/non-xlsx-file.txt", 1)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, `bad xslx file`, w.Body.String())
	})

	name = "post completly correct table with sellerId 1"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/testtables/01_correct.xlsx", 1)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, `{"added":20,"updated":0,"deleted":0,"errors":[]}`, w.Body.String())
	})

	name = "post completly correct table with sellerId 2 also"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/testtables/02_correct.xlsx", 2)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, `{"added":20,"updated":0,"deleted":0,"errors":[]}`, w.Body.String())
	})

	name = "post table with errors with sellerId 3"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/xlsxparser/test/example_with_errors.xlsx", 3)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, respBodyWithErrors, w.Body.String())
	})

	name = "post table with delete offerIDs 11..20 with sellerId 1"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/testtables/03_correct_delete.xlsx", 1)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, `{"added":0,"updated":0,"deleted":10,"errors":[]}`, w.Body.String())
	})

	name = "post table with update offerIDs 1..4,19,20 with sellerId 2"
	t.Run(name, func(t *testing.T) {
		log.Println(name)

		w, r := prepreWR(t, "/testtables/04_correct_update.xlsx", 2)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, `{"added":0,"updated":6,"deleted":0,"errors":[]}`, w.Body.String())
	})

}

func prepreWR(t *testing.T, path string, sellerId uint64) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()

	buf, err := json.Marshal(postTableJson{
		TableURL: tableReqHostPort + path,
		SellerId: sellerId,
	})
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	bodyReader := bytes.NewBuffer(buf)

	testRequest := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	return w, testRequest
}
