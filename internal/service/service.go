package service

import "github.com/hablof/product-registration/internal/models"

type Service struct {
}

type Repo interface {
}

type RequestFilter struct {
	SellerId  uint64
	OfferId   uint64
	Substring string
}

type UpdateStatus struct {
	Added   uint64
	Updated uint64
	Errors  []error
}

func (s *Service) UpdateProducts(sellerId uint64, productsInfo []models.ProductUpdate) (UpdateStatus, error) {

	panic("unimplemented")
}

func (s *Service) ProductsByFilter(filter RequestFilter) ([]models.Product, error) {

	panic("unimplemented")
}
