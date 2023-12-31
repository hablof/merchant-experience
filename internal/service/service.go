package service

import (
	"errors"
	"log"
	"sort"

	"github.com/hablof/merchant-experience/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	s := Service{
		repo: r,
	}
	return &s
}

type Repository interface {
	// SellerProducts(sellerId uint64) ([]models.Product, error)
	SellerProductIDs(sellerId uint64) ([]uint64, error)

	// AddProducts(sellerId uint64, products []models.Product) error

	// единый метод для обеспечения транзакционности внутри репо
	ManageProducts(
		sellerId uint64,
		productsToAdd []models.Product,
		productsToDelete []models.Product,
		productsToUpdate []models.Product,
	) (uint64, error)

	ProductsByFilter(filter RequestFilter) ([]models.Product, error)
}

// type ManageProductsError struct {
// 	Errors error
// }

type RequestFilter struct {
	SellerIDs []uint64
	OfferIDs  []uint64
	Substring string
}

type UpdateResults struct {
	Added   uint64  `json:"added"`
	Updated uint64  `json:"updated"`
	Deleted uint64  `json:"deleted"`
	Errors  []error `json:"errors"`
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
	toAdd := make([]models.Product, 0)
	toUpd := make([]models.Product, 0)
	toDel := make([]models.Product, 0)

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
	validToDel := make([]models.Product, 0, len(toDel))
	validationErrs := make([]error, 0)

	for _, product := range toAdd {
		if err := product.Validate(); err != nil {
			validationErrs = append(validationErrs, err)
		} else {
			validToAdd = append(validToAdd, product)
		}
	}
	for _, product := range toUpd {
		if err := product.Validate(); err != nil {
			validationErrs = append(validationErrs, err)
		} else {
			validToUpd = append(validToUpd, product)
		}
	}
	validToDel = append(validToDel, toDel...) // не знаю как на тестах положительно сравнить одинаково наполненные слайсы с разной capacity

	if len(validToAdd) == 0 && len(validToDel) == 0 && len(validToUpd) == 0 {
		ur := UpdateResults{}
		ur.Errors = append(ur.Errors, validationErrs...)
		return ur, nil
	}

	actualDeleted, err := s.repo.ManageProducts(sellerId, validToAdd, validToDel, validToUpd)
	if err != nil {
		log.Println(err)
		return UpdateResults{}, errors.New("repo err")
	}

	totalErrors := make([]error, 0, len(validationErrs))
	totalErrors = append(totalErrors, validationErrs...)

	return UpdateResults{
		Added:   uint64(len(validToAdd)),
		Updated: uint64(len(validToUpd)),
		Deleted: actualDeleted,
		Errors:  totalErrors,
	}, nil
}

func contains(slice []uint64, elem uint64) bool {
	if len(slice) == 0 {
		return false
	}

	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= elem })
	if i >= len(slice) || i < 0 {
		return false
	}

	return slice[i] == elem
}

func (s *Service) ProductsByFilter(filter RequestFilter) ([]models.Product, error) {

	// filter.Substring = strings.TrimSpace(filter.Substring)
	products, err := s.repo.ProductsByFilter(filter)
	if err != nil {
		log.Println(err)
		return nil, errors.New("repo err")
	}

	return products, nil
}
