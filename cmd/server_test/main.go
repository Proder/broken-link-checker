package main

import (
	"broken-link-checker/internal/config"
	"broken-link-checker/internal/delivery/http_test"
	"log"
)

func main() {
	// Get the server settings
	cnf, err := config.Get()
	if err != nil {
		log.Fatal("config.Get failed: ", err)
	}

	// Start the server for testing
	if err := http_test.StartServer(&cnf.ServerTest); err != nil {
		log.Fatal("http_test.StartServer failed: ", err)
	}
}
