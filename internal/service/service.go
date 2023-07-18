package service

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/hablof/product-registration/internal/models"
)

type Service struct {
	repo Repo
}

type Repo interface {
	SellerProducts(sellerId uint64) ([]models.Product, error)
	SellerProductIDs(sellerId uint64) ([]uint64, error)

	// AddProducts(sellerId uint64, products []models.Product) error

	// единый метод для обеспечения транзакционности внутри репо
	ManageProducts(
		sellerId uint64,
		productsToAdd []models.Product,
		productsToDelete []models.Product,
		productsToUpdate []models.Product,
	) (UpdateResults, error)
}

type ManageProductsError struct {
	Errors error
}

type RequestFilter struct {
	SellerId  uint64
	OfferId   uint64
	Substring string
}

type UpdateResults struct {
	Added   uint64
	Updated uint64
	Deleted uint64
	Errors  []error
}

func (s *Service) UpdateProducts(sellerId uint64, productUpdates []models.ProductUpdate) (UpdateResults, error) {

	if len(productUpdates) == 0 {
		return UpdateResults{}, errors.New("empty request")
	}

	sellerProductIDs, err := s.repo.SellerProductIDs(sellerId)
	if err != nil {
		log.Println(err)
		return UpdateResults{}, errors.New("repo err")
	}

	if !sort.SliceIsSorted(sellerProductIDs, func(i, j int) bool { return sellerProductIDs[i] < sellerProductIDs[j] }) {
		sort.Slice(sellerProductIDs, func(i, j int) bool { return sellerProductIDs[i] < sellerProductIDs[j] })
	}
	// validatedUpdates := make([]models.ProductUpdate, 0, len(productUpdates))
	var (
		toAdd []models.Product
		toDel []models.Product
		toUpd []models.Product
	)

	for _, upd := range productUpdates {
		switch {
		case !upd.Available:
			toDel = append(toDel, upd.Product)

		case contains(sellerProductIDs, upd.Product.OfferId):
			toUpd = append(toUpd, upd.Product)

		default:
			toAdd = append(toAdd, upd.Product)
		}
	}

	validToAdd := make([]models.Product, 0, len(toAdd))
	validToUpd := make([]models.Product, 0, len(toUpd))
	validationErrs := make([]error, 0)

	for _, product := range toAdd {
		if err := product.Validate(); err != nil {
			validationErrs = append(validationErrs, fmt.Errorf("product invalid: id=%d, err=%v", product.OfferId, err))
		} else {
			validToAdd = append(validToAdd, product)
		}
	}
	for _, product := range toUpd {
		if err := product.Validate(); err != nil {
			validationErrs = append(validationErrs, fmt.Errorf("product invalid: id=%d, err=%v", product.OfferId, err))
		} else {
			validToUpd = append(validToUpd, product)
		}
	}

	repoUpdateResult, err := s.repo.ManageProducts(sellerId, validToAdd, toDel, validToUpd)
	if err != nil {
		log.Println(err)
		return UpdateResults{}, errors.New("repo err")
	}

	totalErrors := make([]error, 0, len(validationErrs)+len(repoUpdateResult.Errors))
	for _, err := range validationErrs {
		totalErrors = append(totalErrors, err)
	}
	for _, err := range repoUpdateResult.Errors {
		totalErrors = append(totalErrors, err)
	}

	return UpdateResults{
		Added:   repoUpdateResult.Added,
		Updated: repoUpdateResult.Updated,
		Deleted: repoUpdateResult.Deleted,
		Errors:  totalErrors,
	}, nil
}

func contains(s []uint64, elem uint64) bool {
	i := sort.Search(len(s), func(i int) bool { return s[i] >= elem })
	return s[i] == elem
}

func (s *Service) ProductsByFilter(filter RequestFilter) ([]models.Product, error) {

	panic("unimplemented")
}
