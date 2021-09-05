package main

import (
	"broken-link-checker/internal/config"
	"broken-link-checker/internal/delivery/http"
	"log"
)

func main() {
	// Get the server settings
	cnf := config.Get()

	// Start the server
	if err := http.StartServer(&cnf.Server); err != nil {
		log.Fatal("http.StartServer failed: ", err)
	}
}
