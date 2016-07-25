package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"

	"github.com/erkki/hello-nomad/handlers"
)

func main() {
	log.Println("Starting hello-nomad...")

	httpPort := os.Getenv("NOMAD_PORT_http")
	if httpPort == "" {
		log.Fatal("NOMAD_ADDR_http must be set and non-empty")
	}
	log.Printf("HTTP service listening on port %s", httpPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)

	httpServer := manners.NewServer()
	httpServer.Addr = fmt.Sprintf(":%s", httpPort)
	httpServer.Handler = handlers.LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
