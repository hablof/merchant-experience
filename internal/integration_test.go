package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hablof/product-registration/internal/config"
	"github.com/hablof/product-registration/internal/database"
	"github.com/hablof/product-registration/internal/gateway"
	"github.com/hablof/product-registration/internal/pkg/testfileserver"
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
	sellerN2Updated    = `[{"sellerId":2,"offerId":1,"name":"head_updated","price":1000,"quantity":1000},{"sellerId":2,"offerId":2,"name":"body_updated","price":1000,"quantity":1000},{"sellerId":2,"offerId":3,"name":"name1_3_updated","price":1000,"quantity":1000},{"sellerId":2,"offerId":4,"name":"bigchangus_updated","price":1000,"quantity":1000},{"sellerId":2,"offerId":19,"name":"name15_10_updated","price":1000,"quantity":1000},{"sellerId":2,"offerId":20,"name":"subtitles_updated","price":1000,"quantity":1000}]`
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
		Handler:     testfileserver.NewTestServerHandler(),
		ReadTimeout: 1 * time.Second,
	}
	go func(testServer *http.Server) {
		if err := testServer.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			assert.FailNow(t, err.Error())
		}
	}(fileServer)
	defer fileServer.Close()

	cfg := config.Config{
		Server:     config.Server{Timeout: 5},
		Database:   config.Database{HostLocal: "localhost", Port: "5432", User: "postgres", Password: "1234", DBName: "integration_testing"},
		Repository: config.Repository{Timeout: 5},
		Gateway:    config.Gateway{Timeout: 5},
	}

	db, err := database.NewPostgres(cfg, false)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "no database connection")
	}

	r := repository.NewRepository(db, cfg)
	s := service.NewService(r)
	g := gateway.NewGateway(cfg)
	p := xlsxparser.NewParser()
	handler := router.NewRouter(s, g, p)

	databaseSetup(t, db)
	defer databaseTeardown(t, db)

	// ЗАПИСЬ
	testsPostTable := []struct {
		name        string
		pathToTable string
		sellerId    uint64

		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "post table with duplicates",
			pathToTable:    "/xlsxparser/test/example_duplicates.xlsx",
			sellerId:       42,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "xslx file has duplicates",
		},
		{
			name:           "post empty table",
			pathToTable:    "/xlsxparser/test/example_empty.xlsx",
			sellerId:       42,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "bad xslx file",
		},
		{
			name:           "post table with invalid offer_id column",
			pathToTable:    "/xlsxparser/test/example_with_invalid_offer_id_col.xlsx",
			sellerId:       42,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "offer_id column has invalid value(s)",
		},
		{
			name:           "post txt file",
			pathToTable:    "/testtables/non-xlsx-file.txt",
			sellerId:       42,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "bad xslx file",
		},
		{
			name:           "post completly correct table with sellerId 1",
			pathToTable:    "/testtables/01_correct.xlsx",
			sellerId:       1,
			wantStatusCode: http.StatusOK,
			wantBody:       `{"added":20,"updated":0,"deleted":0,"errors":[]}`,
		},
		{
			name:           "post completly correct table with sellerId 2 also",
			pathToTable:    "/testtables/02_correct.xlsx",
			sellerId:       2,
			wantStatusCode: http.StatusOK,
			wantBody:       `{"added":20,"updated":0,"deleted":0,"errors":[]}`,
		},
		{
			name:           "post table with errors with sellerId 3",
			pathToTable:    "/xlsxparser/test/example_with_errors.xlsx",
			sellerId:       3,
			wantStatusCode: http.StatusOK,
			wantBody:       respBodyWithErrors,
		},
		{
			name:           "post table with delete offerIDs 11..20 with sellerId 1",
			pathToTable:    "/testtables/03_correct_delete.xlsx",
			sellerId:       1,
			wantStatusCode: http.StatusOK,
			wantBody:       `{"added":0,"updated":0,"deleted":10,"errors":[]}`,
		},
		{
			name:           "post table with update offerIDs 1..4,19,20 with sellerId 2",
			pathToTable:    "/testtables/04_correct_update.xlsx",
			sellerId:       2,
			wantStatusCode: http.StatusOK,
			wantBody:       `{"added":0,"updated":6,"deleted":0,"errors":[]}`,
		},
	}

	for _, tt := range testsPostTable {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)

			w, r := preparePostWR(t, tt.pathToTable, tt.sellerId)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatusCode, w.Result().StatusCode)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}

	// ЧТЕНИЕ
	testsGet := []struct {
		name       string
		pSellerIDs string
		pOfferIDs  string
		pSubstring string

		wantStatusCode int
		wantBody       string
	}{
		{
			name:           "updated seller's #2",
			pSellerIDs:     "2",
			pOfferIDs:      "",
			pSubstring:     "updated",
			wantStatusCode: 200,
			wantBody:       sellerN2Updated,
		},
		{
			name:           "",
			pSellerIDs:     "",
			pOfferIDs:      "1,2",
			pSubstring:     "bo",
			wantStatusCode: 200,
			wantBody:       `[{"sellerId":1,"offerId":2,"name":"body_1","price":20,"quantity":0},{"sellerId":3,"offerId":2,"name":"body","price":20,"quantity":0},{"sellerId":2,"offerId":2,"name":"body_updated","price":1000,"quantity":1000}]`,
		},
	}

	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)

			w, r := prepareGetWR(t, tt.pSellerIDs, tt.pOfferIDs, tt.pSubstring)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatusCode, w.Result().StatusCode)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func preparePostWR(t *testing.T, path string, sellerId uint64) (*httptest.ResponseRecorder, *http.Request) {
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

func prepareGetWR(t *testing.T, pSellerIDs string, pOfferIDs string, pSubstring string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()

	paramVals := url.Values{}
	if pSellerIDs != "" {
		paramVals.Add("seller_id", pSellerIDs)
	}
	if pOfferIDs != "" {
		paramVals.Add("offer_id", pOfferIDs)
	}
	if pSubstring != "" {
		paramVals.Add("substring", pSubstring)
	}

	testRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	testRequest.URL.RawQuery = paramVals.Encode()

	return w, testRequest
}
