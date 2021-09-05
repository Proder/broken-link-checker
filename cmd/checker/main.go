package main

import (
	"broken-link-checker/internal/config"
	"broken-link-checker/internal/delivery/http"
	"log"
)

func main() {
	// Get the server settings
	cnf, err := config.Get()
	if err != nil {
		log.Fatal("config.Get failed: ", err)
	}

	// Start the server
	if err := http.StartServer(&cnf.Server); err != nil {
		log.Fatal("http.StartServer failed: ", err)
	}
}
