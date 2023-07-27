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

	r := repository.NewRepository(db, cfg)
	s := service.NewService(r)
	g := gateway.NewGateway(cfg)
	p := xlsxparser.NewParser()
	handler := router.NewRouter(s, g, p)

	server := &http.Server{
		Addr:        cfg.Server.Port,
		Handler:     handler,
		ReadTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	log.Println("starting server...")

	go func(server *http.Server) {
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
