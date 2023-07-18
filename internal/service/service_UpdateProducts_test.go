package service

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/hablof/product-registration/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProducts(t *testing.T) {
	testCases := []struct {
		name           string
		sellerId       uint64
		productUpdates []models.ProductUpdate

		mSellerProductIDsExpects    uint64
		mSellerProductIDsReturns    []uint64
		mSellerProductIDsReturnsErr error
		mSellerProductIDsBehavior   func(rMock *RepoMock, expectedInput uint64, returns []uint64, returnsErr error)

		mManageProductsExpectedToAdd []models.Product
		mManageProductsExpectedToUpd []models.Product
		mManageProductsExpectedToDel []models.Product
		mManageProductsReturns       UpdateResults
		mManageProductsReturnsErr    error
		mManageProductsBehavior      func(rMock *RepoMock, expSellerId uint64, expToAdd []models.Product, expToUpd []models.Product, expToDel []models.Product, returns []uint64, returnsErr error)

		shouldReturn UpdateResults
		returnsError error
	}{
		{
			name:           "пустой запрос",
			sellerId:       0,
			productUpdates: []models.ProductUpdate{},

			mSellerProductIDsBehavior: func(rMock *RepoMock, expectedInput uint64, returns []uint64, returnsErr error) {},
			mManageProductsBehavior: func(rMock *RepoMock, expSellerId uint64, expToAdd, expToUpd, expToDel []models.Product, returns []uint64, returnsErr error) {
			},

			returnsError: errors.New("empty request"),
		},
		{
			name:     "валидный запрос: только добавление",
			sellerId: 1,
			productUpdates: []models.ProductUpdate{
				{
					Product: models.Product{
						OfferId:  1,
						Name:     "test1",
						Price:    100,
						Quantity: 5,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  2,
						Name:     "test2",
						Price:    1000,
						Quantity: 50,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  3,
						Name:     "test3",
						Price:    10,
						Quantity: 20,
					},
					Available: true,
				},
			},

			mSellerProductIDsReturns:    []uint64{},
			mSellerProductIDsReturnsErr: nil,
			mSellerProductIDsBehavior: func(rMock *RepoMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},

			mManageProductsReturns: UpdateResults{
				Added:   3,
				Updated: 0,
				Deleted: 0,
				Errors:  nil,
			},
			mManageProductsReturnsErr: nil,
			mManageProductsBehavior: func(rMock *RepoMock, expSellerId uint64, expectedToAdd, expectedToUpd, expectedToDel []models.Product, returns []uint64, returnsErr error) {
			},

			shouldReturn: UpdateResults{
				Added:   3,
				Updated: 0,
				Deleted: 0,
				Errors:  nil,
			},
			returnsError: nil,
		},

		{
			name:     "валидный запрос: только обновление",
			sellerId: 1,
			productUpdates: []models.ProductUpdate{
				{
					Product: models.Product{
						OfferId:  1,
						Name:     "test11",
						Price:    100,
						Quantity: 5,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  2,
						Name:     "test22",
						Price:    1000,
						Quantity: 50,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  3,
						Name:     "test33",
						Price:    10,
						Quantity: 20,
					},
					Available: true,
				},
			},

			mSellerProductIDsReturns:    []uint64{1, 2, 3},
			mSellerProductIDsReturnsErr: nil,
			mManageProductsReturns: UpdateResults{
				Added:   0,
				Updated: 3,
				Deleted: 0,
				Errors:  nil,
			},
			mManageProductsReturnsErr: nil,

			shouldReturn: UpdateResults{
				Added:   0,
				Updated: 3,
				Deleted: 0,
				Errors:  nil,
			},
			returnsError: nil,
		},
		{
			name:     "валидный запрос: только удаление",
			sellerId: 1,
			productUpdates: []models.ProductUpdate{
				{
					Product: models.Product{
						OfferId:  1,
						Name:     "test11",
						Price:    100,
						Quantity: 5,
					},
					Available: false,
				},
				{
					Product: models.Product{
						OfferId:  2,
						Name:     "test22",
						Price:    1000,
						Quantity: 50,
					},
					Available: false,
				},
				{
					Product: models.Product{
						OfferId:  3,
						Name:     "test33",
						Price:    10,
						Quantity: 20,
					},
					Available: false,
				},
			},

			mSellerProductIDsReturns:    []uint64{1, 2, 3},
			mSellerProductIDsReturnsErr: nil,
			mManageProductsReturns: UpdateResults{
				Added:   0,
				Updated: 0,
				Deleted: 3,
				Errors:  nil,
			},
			mManageProductsReturnsErr: nil,

			shouldReturn: UpdateResults{
				Added:   0,
				Updated: 0,
				Deleted: 3,
				Errors:  nil,
			},
			returnsError: nil,
		},

		{
			name:                        "",
			sellerId:                    0,
			productUpdates:              []models.ProductUpdate{},
			mSellerProductIDsReturns:    []uint64{},
			mSellerProductIDsReturnsErr: nil,
			mManageProductsReturns:      UpdateResults{},
			mManageProductsReturnsErr:   nil,
			shouldReturn:                UpdateResults{},
			returnsError:                nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			rMock := NewRepoMock(mc)
			tc.mSellerProductIDsBehavior(rMock, tc.mSellerProductIDsExpects, tc.mSellerProductIDsReturns, tc.mSellerProductIDsReturnsErr)
			tc.mManageProductsBehavior(rMock, tc.sellerId, tc.mManageProductsExpectedToAdd, tc.mManageProductsExpectedToUpd, tc.mManageProductsExpectedToDel, tc.mSellerProductIDsReturns, tc.mManageProductsReturnsErr)
			s := Service{
				repo: rMock,
			}
			actualResult, actualErr := s.UpdateProducts(tc.sellerId, tc.productUpdates)
			assert.Equal(t, tc.shouldReturn, actualResult, "")
			assert.ElementsMatch(t, tc.mManageProductsReturns.Errors, actualResult.Errors)
			assert.ErrorIs(t, actualErr, tc.returnsError)
		})
	}
}
