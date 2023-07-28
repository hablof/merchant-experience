package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hablof/product-registration/internal/config"
	"github.com/hablof/product-registration/internal/database"
	"github.com/hablof/product-registration/internal/gateway"
	"github.com/hablof/product-registration/internal/repository"
	"github.com/hablof/product-registration/internal/router"
	"github.com/hablof/product-registration/internal/service"
	"github.com/hablof/product-registration/internal/xlsxparser"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.ReadConfigYml("config.yml")
	if err != nil {
		log.Println(err)
		return
	}

	inDocker := false
	if os.Getenv("CONTAINER") != "" {
		inDocker = true
	}

	db, err := database.NewPostgres(cfg, inDocker)
	if err != nil {
		log.Printf("no database connection: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("failed to set goose dialect %v", err)
		return
	}

	currentDBVersion, err := goose.EnsureDBVersion(db.DB)
	if err != nil {
		log.Printf("failed to set goose dialect %v", err)
		return
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {

		if err := goose.DownTo(db.DB, "migrations", currentDBVersion); err != nil {
			log.Printf("failed to DOWN migrations: %v", err)
		}

		log.Printf("failed to UP migrations %v", err)
		return
	}

	r := repository.NewRepository(db, cfg)
	s := service.NewService(r)
	g := gateway.NewGateway(cfg)
	p := xlsxparser.NewParser()
	handler := router.NewRouter(s, g, p)

	server := &http.Server{
		Addr:        ":" + cfg.Server.Port,
		Handler:     handler,
		ReadTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	go func(server *http.Server) {
		log.Printf("starting server on %s ...", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Println(err)
		}
	}(server)

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, os.Interrupt, syscall.SIGTERM)

	<-terminationChannel
	log.Println("terminating server...")
	server.Close()
}
