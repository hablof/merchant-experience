package service

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/hablof/merchant-experience/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_contains(t *testing.T) {

	tests := []struct {
		name  string
		slice []uint64
		elem  uint64
		want  bool
	}{
		{
			name:  "первый элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5},
			elem:  0,
			want:  true,
		},
		{
			name:  "последний элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  6,
			want:  true,
		},
		{
			name:  "центральный элемент",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  3,
			want:  true,
		},
		{
			name:  "отсутствующий элемент в середине",
			slice: []uint64{0, 1, 2, 4, 5, 6},
			elem:  3,
			want:  false,
		},
		{
			name:  "отсутствующий элемент больше большего",
			slice: []uint64{0, 1, 2, 3, 4, 5, 6},
			elem:  7,
			want:  false,
		},
		{
			name:  "отсутствующий элемент меньше меньшего",
			slice: []uint64{1, 2, 3, 4, 5, 6},
			elem:  0,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.slice, tt.elem)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateProducts(t *testing.T) {
	testCases := []struct {
		name           string
		sellerId       uint64
		productUpdates []models.ProductUpdate

		mSellerProductIDs_Expects    uint64
		mSellerProductIDs_Returns    []uint64
		mSellerProductIDs_ReturnsErr error
		mSellerProductIDs_Behavior   func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error)

		mManageProducts_ExpectedToAdd  []models.Product
		mManageProducts_ExpectedToUpd  []models.Product
		mManageProducts_ExpectedToDel  []models.Product
		mManageProducts_ReturnsDeleted uint64
		mManageProducts_ReturnsErr     error
		mManageProducts_Behavior       func(rMock *RepositoryMock, expSellerId uint64, expToAdd []models.Product, expToUpd []models.Product, expToDel []models.Product, returns uint64, returnsErr error)

		shouldReturn UpdateResults
		returnsError error
	}{
		{
			name:           "пустой запрос",
			sellerId:       0,
			productUpdates: []models.ProductUpdate{},

			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {},
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expToAdd, expToUpd, expToDel []models.Product, returns uint64, returnsErr error) {
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

			mSellerProductIDs_Expects:    1,
			mSellerProductIDs_Returns:    []uint64{},
			mSellerProductIDs_ReturnsErr: nil,
			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},

			mManageProducts_ExpectedToAdd: []models.Product{
				{
					OfferId:  1,
					Name:     "test1",
					Price:    100,
					Quantity: 5,
				},
				{
					OfferId:  2,
					Name:     "test2",
					Price:    1000,
					Quantity: 50,
				},
				{
					OfferId:  3,
					Name:     "test3",
					Price:    10,
					Quantity: 20,
				},
			},
			mManageProducts_ExpectedToUpd:  []models.Product{},
			mManageProducts_ExpectedToDel:  []models.Product{},
			mManageProducts_ReturnsDeleted: uint64(0),
			mManageProducts_ReturnsErr:     nil,
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expectedToAdd, expectedToUpd, expectedToDel []models.Product, returns uint64, returnsErr error) {
				rMock.ManageProductsMock.Expect(1, expectedToAdd, expectedToDel, expectedToUpd).Return(returns, returnsErr)
			},

			shouldReturn: UpdateResults{
				Added:   3,
				Updated: 0,
				Deleted: 0,
				Errors:  []error{},
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

			mSellerProductIDs_Expects:    1,
			mSellerProductIDs_Returns:    []uint64{1, 2, 3},
			mSellerProductIDs_ReturnsErr: nil,
			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},

			mManageProducts_ExpectedToAdd: []models.Product{},
			mManageProducts_ExpectedToUpd: []models.Product{
				{
					OfferId:  1,
					Name:     "test11",
					Price:    100,
					Quantity: 5,
				},
				{
					OfferId:  2,
					Name:     "test22",
					Price:    1000,
					Quantity: 50,
				},
				{
					OfferId:  3,
					Name:     "test33",
					Price:    10,
					Quantity: 20,
				},
			},
			mManageProducts_ExpectedToDel:  []models.Product{},
			mManageProducts_ReturnsDeleted: uint64(0),
			mManageProducts_ReturnsErr:     nil,

			shouldReturn: UpdateResults{
				Added:   0,
				Updated: 3,
				Deleted: 0,
				Errors:  []error{},
			},
			returnsError: nil,
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expToAdd, expToUpd, expToDel []models.Product, returns uint64, returnsErr error) {
				rMock.ManageProductsMock.Expect(expSellerId, expToAdd, expToDel, expToUpd).Return(returns, returnsErr)
			},
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

			mSellerProductIDs_Returns:    []uint64{1, 2, 3},
			mSellerProductIDs_ReturnsErr: nil,
			mSellerProductIDs_Expects:    1,
			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},

			mManageProducts_ExpectedToAdd: []models.Product{},
			mManageProducts_ExpectedToUpd: []models.Product{},
			mManageProducts_ExpectedToDel: []models.Product{
				{
					OfferId:  1,
					Name:     "test11",
					Price:    100,
					Quantity: 5,
				},
				{
					OfferId:  2,
					Name:     "test22",
					Price:    1000,
					Quantity: 50,
				},
				{
					OfferId:  3,
					Name:     "test33",
					Price:    10,
					Quantity: 20,
				},
			},
			mManageProducts_ReturnsDeleted: uint64(3),
			mManageProducts_ReturnsErr:     nil,
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expToAdd, expToUpd, expToDel []models.Product, returns uint64, returnsErr error) {
				rMock.ManageProductsMock.Expect(expSellerId, expToAdd, expToDel, expToUpd).Return(returns, returnsErr)
			},

			shouldReturn: UpdateResults{
				Added:   0,
				Updated: 0,
				Deleted: 3,
				Errors:  []error{},
			},
			returnsError: nil,
		},
		{
			name:     "валидный запрос: два добавления, два удаления, два обновления",
			sellerId: 42,
			productUpdates: []models.ProductUpdate{
				{
					Product: models.Product{
						OfferId:  1,
						Name:     "первый кандидат на удаление",
						Price:    15,
						Quantity: 15,
					},
					Available: false,
				},
				{
					Product: models.Product{
						OfferId:  2,
						Name:     "первый кандидат на добавление",
						Price:    125,
						Quantity: 10,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  3,
						Name:     "первый кандидат на обновление",
						Price:    4990,
						Quantity: 1,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  4,
						Name:     "второй кандидат на удаление",
						Price:    200,
						Quantity: 5,
					},
					Available: false,
				},
				{
					Product: models.Product{
						OfferId:  5,
						Name:     "второй кандидат на добавление",
						Price:    125,
						Quantity: 10,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  6,
						Name:     "второй кандидат на обновление",
						Price:    15000,
						Quantity: 60,
					},
					Available: true,
				},
			},

			mSellerProductIDs_Expects:    42,
			mSellerProductIDs_Returns:    []uint64{1, 3, 4, 6},
			mSellerProductIDs_ReturnsErr: nil,
			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},

			mManageProducts_ExpectedToAdd: []models.Product{
				{
					OfferId:  2,
					Name:     "первый кандидат на добавление",
					Price:    125,
					Quantity: 10,
				},
				{
					OfferId:  5,
					Name:     "второй кандидат на добавление",
					Price:    125,
					Quantity: 10,
				},
			},
			mManageProducts_ExpectedToUpd: []models.Product{
				{
					OfferId:  3,
					Name:     "первый кандидат на обновление",
					Price:    4990,
					Quantity: 1,
				},
				{
					OfferId:  6,
					Name:     "второй кандидат на обновление",
					Price:    15000,
					Quantity: 60,
				},
			},
			mManageProducts_ExpectedToDel: []models.Product{
				{
					OfferId:  1,
					Name:     "первый кандидат на удаление",
					Price:    15,
					Quantity: 15,
				},
				{
					OfferId:  4,
					Name:     "второй кандидат на удаление",
					Price:    200,
					Quantity: 5,
				},
			},
			mManageProducts_ReturnsDeleted: uint64(2),
			mManageProducts_ReturnsErr:     nil,
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expToAdd []models.Product, expToUpd []models.Product, expToDel []models.Product, returns uint64, returnsErr error) {
				rMock.ManageProductsMock.Expect(expSellerId, expToAdd, expToDel, expToUpd).Return(returns, returnsErr)
			},
			shouldReturn: UpdateResults{
				Added:   2,
				Updated: 2,
				Deleted: 2,
				Errors:  []error{},
			},
			returnsError: nil,
		},
		{
			name:     "добавляем и обновляем только длинные (невалидные) названия",
			sellerId: 70,
			productUpdates: []models.ProductUpdate{
				{
					Product: models.Product{
						OfferId:  15,
						Name:     "ну очень очень очень очень очень очень очень очень очень очень очень очень очень очень очень длинное название",
						Price:    0,
						Quantity: 0,
					},
					Available: true,
				},
				{
					Product: models.Product{
						OfferId:  16,
						Name:     "ну очень очень очень очень очень очень очень очень очень очень очень очень очень очень очень длинное название",
						Price:    0,
						Quantity: 0,
					},
					Available: true,
				},
			},
			mSellerProductIDs_Expects:    70,
			mSellerProductIDs_Returns:    []uint64{15},
			mSellerProductIDs_ReturnsErr: nil,
			mSellerProductIDs_Behavior: func(rMock *RepositoryMock, expectedInput uint64, returns []uint64, returnsErr error) {
				rMock.SellerProductIDsMock.Expect(expectedInput).Return(returns, returnsErr)
			},
			mManageProducts_ExpectedToAdd:  []models.Product{},
			mManageProducts_ExpectedToUpd:  []models.Product{},
			mManageProducts_ExpectedToDel:  []models.Product{},
			mManageProducts_ReturnsDeleted: uint64(0),
			mManageProducts_ReturnsErr:     nil,
			mManageProducts_Behavior: func(rMock *RepositoryMock, expSellerId uint64, expToAdd []models.Product, expToUpd []models.Product, expToDel []models.Product, returns uint64, returnsErr error) {
			},
			shouldReturn: UpdateResults{
				Added:   0,
				Updated: 0,
				Deleted: 0,
				Errors: []error{
					models.ErrProductValidation{
						OfferId: 16,
						Field:   "name",
						ErrMsg:  models.MsgTooLongName,
					},
					models.ErrProductValidation{
						OfferId: 15,
						Field:   "name",
						ErrMsg:  models.MsgTooLongName,
					},
				},
			},
			returnsError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			rMock := NewRepositoryMock(mc)
			tc.mSellerProductIDs_Behavior(rMock, tc.mSellerProductIDs_Expects, tc.mSellerProductIDs_Returns, tc.mSellerProductIDs_ReturnsErr)
			tc.mManageProducts_Behavior(rMock, tc.sellerId, tc.mManageProducts_ExpectedToAdd, tc.mManageProducts_ExpectedToUpd, tc.mManageProducts_ExpectedToDel, tc.mManageProducts_ReturnsDeleted, tc.mManageProducts_ReturnsErr)
			s := Service{
				repo: rMock,
			}
			actualResult, actualErr := s.UpdateProducts(tc.sellerId, tc.productUpdates)
			assert.Equal(t, tc.shouldReturn.Added, actualResult.Added, "")
			assert.Equal(t, tc.shouldReturn.Deleted, actualResult.Deleted, "")
			assert.Equal(t, tc.shouldReturn.Updated, actualResult.Updated, "")
			if assert.Equal(t, len(actualResult.Errors), len(tc.shouldReturn.Errors)) {
				for i, elem := range actualResult.Errors {
					assert.Equal(t, tc.shouldReturn.Errors[i], elem)
				}
			}

			// assert.ElementsMatch(t, tc.mManageProducts_Returns.Errors, actualResult.Errors)
			assert.Equal(t, actualErr, tc.returnsError)
		})
	}
}
