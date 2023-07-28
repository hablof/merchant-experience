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

	host := ""
	if inDocker {
		host = cfg.Database.HostDocker
	} else {
		host = cfg.Database.HostLocal
	}

	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName)

	var (
		db  *sqlx.DB
		err error
	)

	// попытки подключения
	for i := 1; i < 11; i++ {
		ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
		defer cf()

		db, err = sqlx.ConnectContext(ctx, "postgres", connectString)

		if err != nil {
			log.Printf("attempt #%d to connect to postgres failed: %v", i, err)
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Println(err)
		return nil, errors.New("failed create connection to postgres")
	}

	log.Println("connected to postgres")

	return db, nil
}
