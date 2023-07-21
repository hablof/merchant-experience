package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres() (*sqlx.DB, error) {
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	db, err := sqlx.ConnectContext(ctx, "postgres", "host=localhost port=5432 user=postgres password=1234 dbname=integration_testing sslmode=disable")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed create connection to postgres")
	}
	// defer db.Close()

	return db, nil
}
