package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hablof/product-registration/internal/pkg/testfileserver"
)

const (
	serverHostPort = ":8015"
)

func main() {
	fileServer := &http.Server{
		Addr:        serverHostPort,
		Handler:     testfileserver.NewTestServerHandler(),
		ReadTimeout: 1 * time.Second,
	}

	log.Println("starting test server...")

	go func(testServer *http.Server) {
		if err := testServer.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			log.Println(err)
		}
	}(fileServer)

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, os.Interrupt, syscall.SIGTERM)

	<-terminationChannel
	log.Println("terminating test server...")
	fileServer.Close()
}
