package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hablof/merchant-experience/internal/models"
	"github.com/hablof/merchant-experience/internal/service"
	xlsxparser "github.com/hablof/merchant-experience/internal/xlsxparser"

	"github.com/julienschmidt/httprouter"
)

const (
	sellerIdParamField  = "seller_id"
	offerIdParamField   = "offer_id"
	substringParamField = "substring"
)

type TableDownloader interface {
	Table(url string) (io.Reader, error)
}

type Service interface {
	ProductsByFilter(filter service.RequestFilter) ([]models.Product, error)
	UpdateProducts(sellerId uint64, productUpdates []models.ProductUpdate) (service.UpdateResults, error)
}

type ExcelParser interface {
	ParseProducts(r io.Reader) (productUpdates []models.ProductUpdate, productErrs []error, methodErr error)
}

type jsonSchema struct {
	TableURL string `json:"tableURL"`
	SellerId uint64 `json:"sellerId"`
}

type Handler struct {
	s  Service
	td TableDownloader
	ep ExcelParser
}

func NewRouter(
	s Service,
	td TableDownloader,
	ep ExcelParser,
) http.Handler {

	h := Handler{
		s:  s,
		td: td,
		ep: ep,
	}

	r := httprouter.New()
	r.GET("/", h.GetProducts)
	r.POST("/", h.PostTableURL)
	r.PanicHandler = h.PanicHanler

	return r
}

func (h *Handler) PanicHanler(w http.ResponseWriter, r *http.Request, _ interface{}) {
	log.Println("panic recovered")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("fatal service error"))
}

func (h *Handler) PostTableURL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	b := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(b); err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unable to read body: " + err.Error())
		fmt.Fprint(w, "unable to read body")

		return
	}

	postStruct := jsonSchema{}
	if err := json.Unmarshal(b, &postStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad json: " + err.Error())
		fmt.Fprint(w, "bad json")

		return
	}

	table, err := h.td.Table(postStruct.TableURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad table url: " + err.Error())
		fmt.Fprint(w, "bad table url")

		return
	}

	productUpdates, productErrs, methodErr := h.ep.ParseProducts(table)
	switch {
	case errors.Is(methodErr, xlsxparser.ErrEmptyDoc),
		errors.Is(methodErr, xlsxparser.ErrEmptySheet),
		errors.Is(methodErr, xlsxparser.ErrFailedToRead):

		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad xslx file")
		fmt.Fprint(w, "bad xslx file")

		return

	case errors.Is(methodErr, xlsxparser.ErrInvalidIDs):
		w.WriteHeader(http.StatusBadRequest)
		log.Println("offer_id column has invalid value(s)")
		fmt.Fprint(w, "offer_id column has invalid value(s)")

		return

	case errors.Is(methodErr, xlsxparser.ErrHasDuplicates):
		w.WriteHeader(http.StatusBadRequest)
		log.Println("xslx file has duplicates")
		fmt.Fprint(w, "xslx file has duplicates")

		return

	case methodErr != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(methodErr.Error())
		fmt.Fprint(w, "parsing error")

		return
	}

	ur, err := h.s.UpdateProducts(postStruct.SellerId, productUpdates)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprint(w, "service error")

		return
	}
	ur.Errors = append(ur.Errors, productErrs...)

	b2, err := json.Marshal(ur)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprint(w, "service error")

		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Type", "charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write(b2)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// fetch url params
	sellerIdParam := r.URL.Query().Get(sellerIdParamField)
	offerIdParam := r.URL.Query().Get(offerIdParamField)
	paramSubstr := r.URL.Query().Get(substringParamField)

	splitedSellerIDsStrs := strings.Split(sellerIdParam, ",")
	splitedOfferIDsStrs := strings.Split(offerIdParam, ",")

	sellerIDs := make([]uint64, 0, len(splitedSellerIDsStrs))
	for _, elem := range splitedSellerIDsStrs {
		u, err := strconv.ParseUint(elem, 10, 64)
		if err != nil {
			sellerIDs = nil
			break
		}

		sellerIDs = append(sellerIDs, u)
	}

	offerIDs := make([]uint64, 0, len(splitedOfferIDsStrs))
	for _, elem := range splitedOfferIDsStrs {
		u, err := strconv.ParseUint(elem, 10, 64)
		if err != nil {
			offerIDs = nil
			break
		}

		offerIDs = append(offerIDs, u)
	}

	rf := service.RequestFilter{
		SellerIDs: sellerIDs,
		OfferIDs:  offerIDs,
		Substring: paramSubstr,
	}
	products, err := h.s.ProductsByFilter(rf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to fetch products: " + err.Error())
		fmt.Fprint(w, "failed to fetch products")

		return
	}

	b, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to marshal products: " + err.Error())
		fmt.Fprint(w, "service error")

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
