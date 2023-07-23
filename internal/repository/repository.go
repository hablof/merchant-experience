package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hablof/product-registration/internal/models"
	"github.com/hablof/product-registration/internal/service"
	"github.com/jmoiron/sqlx"
)

const (
	tableName   = "products"
	sellerIdCol = "seller_id"
	offerIdCol  = "offer_id"
	nameCol     = "name"
	priceCol    = "price"
	quantityCol = "quantity"
)

const (
	defaultLimit = 100
)

var (
	ErrQueryBuilderFailed = errors.New("query builder failed")
	ErrTxFailed           = errors.New("transaction failed")
	ErrQueryExecFailed    = errors.New("failed to execute query")
	ErrEmptyRequest       = errors.New("empty request")
)

type Repository struct {
	db        *sqlx.DB
	initQuery sq.StatementBuilderType
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:        db,
		initQuery: sq.StatementBuilder.PlaceholderFormat(sq.Dollar), // Postgress
	}
}

func (r *Repository) ManageProducts(
	sellerId uint64,
	productsToAdd []models.Product,
	productsToDelete []models.Product,
	productsToUpdate []models.Product,
) error {

	if len(productsToAdd)+len(productsToUpdate)+len(productsToDelete) == 0 {
		return ErrEmptyRequest
	}

	// start transaction
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Println(err)
		return ErrTxFailed
	}
	defer tx.Rollback()

	// insert query
	if len(productsToAdd)+len(productsToUpdate) > 0 {

		// prepare insert query
		insertQuery := r.initQuery.
			Insert(tableName).
			Columns(sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol)

		for _, elem := range productsToAdd {
			insertQuery = insertQuery.Values(sellerId, elem.OfferId, elem.Name, elem.Price, elem.Quantity)
		}
		for _, elem := range productsToUpdate {
			insertQuery = insertQuery.Values(sellerId, elem.OfferId, elem.Name, elem.Price, elem.Quantity)
		}

		insertQuery = insertQuery.Suffix(
			`ON CONFLICT ON CONSTRAINT no_duplicates DO UPDATE SET
			name = EXCLUDED.name,
			price = EXCLUDED.price,
			quantity = EXCLUDED.quantity`,
		)

		insertQueryString, insertQueryArgs, err := insertQuery.ToSql()
		if err != nil {
			log.Println(err)
			return ErrQueryBuilderFailed
		}

		// execute insert
		insertQueryResult, err := tx.ExecContext(ctx, insertQueryString, insertQueryArgs...)
		if err != nil {
			log.Println(err)
			return ErrQueryExecFailed
		}
		rowsAffected, err := insertQueryResult.RowsAffected()
		if err != nil {
			log.Println(err)
			return ErrQueryExecFailed
		}
		if rowsAffected != int64(len(productsToAdd)+len(productsToUpdate)) {
			log.Println("missmatched sum of products to add/update and affected rows")
		}
	}

	// delete query
	if len(productsToDelete) > 0 {
		// prepare delete query
		deleteIDs := make([]uint64, 0, len(productsToDelete))
		for _, elem := range productsToDelete {
			deleteIDs = append(deleteIDs, elem.OfferId)
		}
		deleteQuery := r.initQuery.Delete(tableName).Where(sq.Eq{sellerIdCol: sellerId, offerIdCol: deleteIDs})
		deleteQueryString, deleteQueryArgs, err := deleteQuery.ToSql()
		if err != nil {
			log.Println(err)
			return ErrQueryBuilderFailed
		}
		// execute delete query
		deleteQueryResult, err := tx.ExecContext(ctx, deleteQueryString, deleteQueryArgs...)
		if err != nil {
			log.Println(err)
			return ErrQueryExecFailed
		}
		rowsAffected, err := deleteQueryResult.RowsAffected()
		if err != nil {
			log.Println(err)
			return ErrQueryExecFailed
		}
		if rowsAffected != int64(len(productsToDelete)) {
			log.Println("missmatched sum of products to delete and affected rows")
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return ErrTxFailed
	}

	return nil
}

func (r *Repository) ProductsByFilter(filter service.RequestFilter) ([]models.Product, error) {
	selectQuery := r.initQuery.Select(sellerIdCol, offerIdCol, nameCol, priceCol, quantityCol)

	if len(filter.SellerIDs) > 0 {
		selectQuery = selectQuery.Where(sq.Eq{sellerIdCol: filter.SellerIDs}) // WHERE ... in (...) construction
	}

	if len(filter.OfferIDs) > 0 {
		selectQuery = selectQuery.Where(sq.Eq{offerIdCol: filter.OfferIDs})
	}

	filter.Substring = strings.TrimSpace(filter.Substring)
	if filter.Substring != "" {
		selectQuery = selectQuery.Where(sq.Like{nameCol: filter.Substring})
	}

	selectQueryString, args, err := selectQuery.Limit(defaultLimit).ToSql()
	if err != nil {
		log.Println(err)
		return nil, ErrQueryBuilderFailed
	}

	ctx, cf := context.WithTimeout(context.Background(), 15*time.Second)
	defer cf()

	products := make([]models.Product, 0)
	if err := r.db.SelectContext(ctx, &products, selectQueryString, args...); err != nil {
		log.Println(err)
		return nil, ErrQueryExecFailed
	}

	return products, nil
}
