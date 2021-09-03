package app

import (
	"broken-link-checker/app/config"
	"broken-link-checker/app/internal/delivery/http"
	"broken-link-checker/app/internal/delivery/http_test"
	"fmt"
)

func Run() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server
	err := http.StartServer(&cnf.Server)
	if err != nil {
		return fmt.Errorf("http.StartServer failed: %w", err)
	}

	return nil
}

func RunServerTest() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server for testing
	err := http_test.StartServer(&cnf.ServerTest)
	if err != nil {
		return fmt.Errorf("http_test.StartServer failed: %w", err)
	}

	return nil
}
