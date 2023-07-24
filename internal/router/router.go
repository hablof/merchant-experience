package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hablof/product-registration/internal/models"
	"github.com/hablof/product-registration/internal/service"
	xlsxparser "github.com/hablof/product-registration/internal/xlsxparser"

	"github.com/julienschmidt/httprouter"
)

type TableDownloader interface {
	Table(url string) (io.Reader, error)
}

type Service interface {
	ProductsByFilter(filter service.RequestFilter) ([]models.Product, error)
	UpdateProducts(sellerId uint64, productUpdates []models.ProductUpdate) (service.UpdateResults, error)
}

type jsonSchema struct {
	tableURL string `json:"tableURL"`
	sellerId uint64 `json:"sellerId"`
}

type Handler struct {
	s  Service
	td TableDownloader
}

func NewRouter(s Service, td TableDownloader) http.Handler {
	h := Handler{}
	r := httprouter.New()
	r.GET("/", h.GetProducts)
	r.POST("/", h.PostTableURL)
	return r
}

func (router *Handler) PostTableURL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	b := []byte{}
	if _, err := r.Body.Read(b); err != nil {
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

	table, err := router.td.Table(postStruct.tableURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad table url: " + err.Error())
		fmt.Fprint(w, "bad table url")

		return
	}

	productUpdates, producrErrs, methodErr := xlsxparser.ParseProducts(table)
	switch {
	case errors.Is(methodErr, xlsxparser.ErrEmptyDoc) || errors.Is(methodErr, xlsxparser.ErrEmptyDoc):
		w.WriteHeader(http.StatusBadRequest)
		log.Println("bad xslx file")
		fmt.Fprint(w, "bad xslx file")

		return

	case errors.Is(methodErr, xlsxparser.ErrHasDuplicates):
		w.WriteHeader(http.StatusBadRequest)
		log.Println("xslx file has duplicates")
		fmt.Fprint(w, "xslx file has duplicates")

		return

	case methodErr != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprint(w, "parsing error")

		return
	}

	ur, err := router.s.UpdateProducts(postStruct.sellerId, productUpdates)

}

func (router *Handler) GetProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
