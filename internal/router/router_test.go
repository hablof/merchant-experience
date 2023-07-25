package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/hablof/product-registration/internal/models"
	"github.com/hablof/product-registration/internal/service"
	xlsxparser "github.com/hablof/product-registration/internal/xlsxparser"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetProducts(t *testing.T) {

	tests := []struct {
		name       string
		pSellerIDs string
		pOfferIDs  string
		pSubstring string

		expectedReqFilter service.RequestFilter
		serviceReturns    []models.Product
		serviceReturnsErr error
		serviceBehaviour  func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error)

		wantStatusCode  int
		wantContentBody string
	}{
		{
			name:              "empty request",
			pSellerIDs:        "",
			pOfferIDs:         "",
			pSubstring:        "",
			expectedReqFilter: service.RequestFilter{SellerIDs: nil, OfferIDs: nil, Substring: ""},
			serviceReturns: []models.Product{
				{SellerId: 1, OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{SellerId: 3, OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			serviceReturnsErr: nil,
			serviceBehaviour: func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error) {
				sm.ProductsByFilterMock.Expect(expRF).Return(serviceRet, serviceRetErr)
			},
			wantStatusCode:  200,
			wantContentBody: `[{"sellerId":1,"offerId":1,"name":"name1","price":1,"quantity":1},{"sellerId":2,"offerId":2,"name":"name2","price":2,"quantity":2},{"sellerId":3,"offerId":3,"name":"name3","price":3,"quantity":3}]`,
		},
		{
			name:              "correct params",
			pSellerIDs:        "1,2,3",
			pOfferIDs:         "2,3,4",
			pSubstring:        "name",
			expectedReqFilter: service.RequestFilter{SellerIDs: []uint64{1, 2, 3}, OfferIDs: []uint64{2, 3, 4}, Substring: "name"},
			serviceReturns: []models.Product{
				{SellerId: 1, OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{SellerId: 3, OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			serviceReturnsErr: nil,
			serviceBehaviour: func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error) {
				sm.ProductsByFilterMock.Expect(expRF).Return(serviceRet, serviceRetErr)
			},
			wantStatusCode:  200,
			wantContentBody: `[{"sellerId":1,"offerId":1,"name":"name1","price":1,"quantity":1},{"sellerId":2,"offerId":2,"name":"name2","price":2,"quantity":2},{"sellerId":3,"offerId":3,"name":"name3","price":3,"quantity":3}]`,
		},
		{
			name:              "incorrect params",
			pSellerIDs:        "1,2,3,incorrect",
			pOfferIDs:         "2,incorrect3,4",
			pSubstring:        "name", // name cannot be incorrect
			expectedReqFilter: service.RequestFilter{SellerIDs: nil, OfferIDs: nil, Substring: "name"},
			serviceReturns: []models.Product{
				{SellerId: 1, OfferId: 1, Name: "name1", Price: 1, Quantity: 1},
				{SellerId: 2, OfferId: 2, Name: "name2", Price: 2, Quantity: 2},
				{SellerId: 3, OfferId: 3, Name: "name3", Price: 3, Quantity: 3},
			},
			serviceReturnsErr: nil,
			serviceBehaviour: func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error) {
				sm.ProductsByFilterMock.Expect(expRF).Return(serviceRet, serviceRetErr)
			},
			wantStatusCode:  200,
			wantContentBody: `[{"sellerId":1,"offerId":1,"name":"name1","price":1,"quantity":1},{"sellerId":2,"offerId":2,"name":"name2","price":2,"quantity":2},{"sellerId":3,"offerId":3,"name":"name3","price":3,"quantity":3}]`,
		},
		{
			name:              "service error",
			pSellerIDs:        "",
			pOfferIDs:         "",
			pSubstring:        "name", // name cannot be incorrect
			expectedReqFilter: service.RequestFilter{SellerIDs: nil, OfferIDs: nil, Substring: "name"},
			serviceReturns:    nil,
			serviceReturnsErr: errors.New("repo err"),
			serviceBehaviour: func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error) {
				sm.ProductsByFilterMock.Expect(expRF).Return(serviceRet, serviceRetErr)
			},
			wantStatusCode:  500,
			wantContentBody: `failed to fetch products`,
		},
		{
			name:              "empty response",
			pSellerIDs:        "",
			pOfferIDs:         "",
			pSubstring:        "", // name cannot be incorrect
			expectedReqFilter: service.RequestFilter{SellerIDs: nil, OfferIDs: nil, Substring: ""},
			serviceReturns:    nil,
			serviceReturnsErr: nil,
			serviceBehaviour: func(sm *ServiceMock, expRF service.RequestFilter, serviceRet []models.Product, serviceRetErr error) {
				sm.ProductsByFilterMock.Expect(expRF).Return(serviceRet, serviceRetErr)
			},
			wantStatusCode:  200,
			wantContentBody: `null`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sm := NewServiceMock(t)
			tdm := NewTableDownloaderMock(t)
			epm := NewExcelParserMock(t)
			h := NewRouter(sm, tdm, epm)

			tt.serviceBehaviour(sm, tt.expectedReqFilter, tt.serviceReturns, tt.serviceReturnsErr)

			paramVals := url.Values{}
			if tt.pSellerIDs != "" {
				paramVals.Add("seller_id", tt.pSellerIDs)
			}
			if tt.pOfferIDs != "" {
				paramVals.Add("offer_id", tt.pOfferIDs)
			}
			if tt.pSubstring != "" {
				paramVals.Add("substring", tt.pSubstring)
			}
			params := paramVals.Encode()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/?"+params, bytes.NewBufferString(""))

			h.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatusCode, w.Result().StatusCode, "status code")
			assert.Equal(t, tt.wantContentBody, w.Body.String(), "response body")
		})
	}
}

func TestHandler_PostTableURL(t *testing.T) {

	tests := []struct {
		name        string
		reqUrl      string
		reqSellerID uint64

		tdExpectURL  string
		tdReturns    io.Reader
		tdReturnsErr error
		tdBehaviour  func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdetErr error)

		parserExpectTable  io.Reader
		parserReturns      []models.ProductUpdate
		parserRetValidErrs []error
		parserReturnsErr   error
		parserBehaviour    func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error)

		expectedSellerID  uint64
		expectedUpdates   []models.ProductUpdate
		serviceReturns    service.UpdateResults
		serviceReturnsErr error
		serviceBehaviour  func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error)

		wantStatusCode  int
		wantContentBody string
	}{
		{
			name:         "bad table url",
			reqUrl:       "example.com/table",
			reqSellerID:  42,
			tdExpectURL:  "example.com/table",
			tdReturns:    nil,
			tdReturnsErr: errors.New("some table downloader err"),
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
			},

			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
			},
			wantStatusCode:  400,
			wantContentBody: "bad table url",
		},
		{
			name:         "bad xslx file",
			reqUrl:       "some.url/t",
			reqSellerID:  1,
			tdExpectURL:  "some.url/t",
			tdReturns:    bytes.NewBufferString("table mock"),
			tdReturnsErr: nil,
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserExpectTable:  bytes.NewBufferString("table mock"),
			parserReturns:      nil,
			parserRetValidErrs: nil,
			parserReturnsErr:   xlsxparser.ErrEmptyDoc,
			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
				epm.ParseProductsMock.Expect(parserExpectTable).Return(pReturns, pRetValidErrs, pRetErr)
			},

			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
			},
			wantStatusCode:  400,
			wantContentBody: "bad xslx file",
		},
		{
			name:         "xslx file has duplicates",
			reqUrl:       "some.url/t",
			reqSellerID:  1,
			tdExpectURL:  "some.url/t",
			tdReturns:    bytes.NewBufferString("table mock"),
			tdReturnsErr: nil,
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserExpectTable:  bytes.NewBufferString("table mock"),
			parserReturns:      nil,
			parserRetValidErrs: nil,
			parserReturnsErr:   xlsxparser.ErrHasDuplicates,
			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
				epm.ParseProductsMock.Expect(parserExpectTable).Return(pReturns, pRetValidErrs, pRetErr)
			},

			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
			},
			wantStatusCode:  400,
			wantContentBody: "xslx file has duplicates",
		},
		{
			name:         "unexpected parser error",
			reqUrl:       "some.url/t",
			reqSellerID:  1,
			tdExpectURL:  "some.url/t",
			tdReturns:    bytes.NewBufferString("table mock"),
			tdReturnsErr: nil,
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserExpectTable:  bytes.NewBufferString("table mock"),
			parserReturns:      nil,
			parserRetValidErrs: nil,
			parserReturnsErr:   errors.New("unexpected parser error"),
			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
				epm.ParseProductsMock.Expect(parserExpectTable).Return(pReturns, pRetValidErrs, pRetErr)
			},

			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
			},
			wantStatusCode:  500,
			wantContentBody: "parsing error",
		},
		{
			name:         "service error",
			reqUrl:       "some.url/t",
			reqSellerID:  1,
			tdExpectURL:  "some.url/t",
			tdReturns:    bytes.NewBufferString("table mock"),
			tdReturnsErr: nil,
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserExpectTable:  bytes.NewBufferString("table mock"),
			parserReturns:      []models.ProductUpdate{{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true}, {Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true}},
			parserRetValidErrs: []error{xlsxparser.ErrProductParsing{Row: 3, Field: "name", ErrMsg: models.MsgTooLongName}, xlsxparser.ErrProductParsing{Row: 4, Field: "price", ErrMsg: (&strconv.NumError{Func: "ParseUint", Num: "0-40", Err: strconv.ErrSyntax}).Error()}},
			parserReturnsErr:   nil,
			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
				epm.ParseProductsMock.Expect(parserExpectTable).Return(pReturns, pRetValidErrs, pRetErr)
			},

			expectedSellerID:  1,
			expectedUpdates:   []models.ProductUpdate{{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true}, {Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true}},
			serviceReturns:    service.UpdateResults{},
			serviceReturnsErr: errors.New("repo err"),
			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
				sm.UpdateProductsMock.Expect(expectedSellerID, expectedUpdates).Return(serviceReturns, serviceRetErr)
			},

			wantStatusCode:  500,
			wantContentBody: "service error",
		},
		{
			name:         "correct request",
			reqUrl:       "some.url/t",
			reqSellerID:  1,
			tdExpectURL:  "some.url/t",
			tdReturns:    bytes.NewBufferString("table mock"),
			tdReturnsErr: nil,
			tdBehaviour: func(tdm *TableDownloaderMock, expURL string, tdRet io.Reader, tdRetErr error) {
				tdm.TableMock.Expect(expURL).Return(tdRet, tdRetErr)
			},

			parserExpectTable:  bytes.NewBufferString("table mock"),
			parserReturns:      []models.ProductUpdate{{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true}, {Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true}},
			parserRetValidErrs: []error{xlsxparser.ErrProductParsing{Row: 3, Field: "name", ErrMsg: models.MsgTooLongName}, xlsxparser.ErrProductParsing{Row: 4, Field: "price", ErrMsg: (&strconv.NumError{Func: "ParseUint", Num: "0-40", Err: strconv.ErrSyntax}).Error()}},
			parserReturnsErr:   nil,
			parserBehaviour: func(epm *ExcelParserMock, parserExpectTable io.Reader, pReturns []models.ProductUpdate, pRetValidErrs []error, pRetErr error) {
				epm.ParseProductsMock.Expect(parserExpectTable).Return(pReturns, pRetValidErrs, pRetErr)
			},

			expectedSellerID:  1,
			expectedUpdates:   []models.ProductUpdate{{Product: models.Product{OfferId: 1, Name: "head", Price: 10, Quantity: 1}, Available: true}, {Product: models.Product{OfferId: 2, Name: "body", Price: 20, Quantity: 0}, Available: true}},
			serviceReturns:    service.UpdateResults{Added: 1, Updated: 1, Deleted: 0, Errors: []error{}},
			serviceReturnsErr: nil,
			serviceBehaviour: func(sm *ServiceMock, expectedSellerID uint64, expectedUpdates []models.ProductUpdate, serviceReturns service.UpdateResults, serviceRetErr error) {
				sm.UpdateProductsMock.Expect(expectedSellerID, expectedUpdates).Return(serviceReturns, serviceRetErr)
			},

			wantStatusCode:  200,
			wantContentBody: `{"added":1,"updated":1,"deleted":0,"errors":[{"row":3,"field":"name","errMsg":"too long name"},{"row":4,"field":"price","errMsg":"strconv.ParseUint: parsing \"0-40\": invalid syntax"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)

			sm := NewServiceMock(t)
			tdm := NewTableDownloaderMock(t)
			epm := NewExcelParserMock(t)
			h := NewRouter(sm, tdm, epm)

			tt.tdBehaviour(tdm, tt.tdExpectURL, tt.tdReturns, tt.tdReturnsErr)
			tt.parserBehaviour(epm, tt.parserExpectTable, tt.expectedUpdates, tt.parserRetValidErrs, tt.parserReturnsErr)
			tt.serviceBehaviour(sm, tt.expectedSellerID, tt.expectedUpdates, tt.serviceReturns, tt.serviceReturnsErr)

			body, err := json.Marshal(jsonSchema{
				TableURL: tt.reqUrl,
				SellerId: tt.reqSellerID,
			})
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

			h.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatusCode, w.Result().StatusCode, "status code")
			assert.Equal(t, tt.wantContentBody, w.Body.String(), "response body")
		})
	}
}
