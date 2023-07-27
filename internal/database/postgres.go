package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hablof/product-registration/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.Config, inDocker bool) (*sqlx.DB, error) {
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	host := ""
	if inDocker {
		host = cfg.Database.HostDocker
	} else {
		host = cfg.Database.HostLocal
	}
	db, err := sqlx.ConnectContext(ctx, "postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
		),
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed create connection to postgres")
	}
	// defer db.Close()

	return db, nil
}
