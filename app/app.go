package app

import (
	"broken-link-checker/app/config"
	http2 "broken-link-checker/app/internal/delivery/http"
	http_test2 "broken-link-checker/app/internal/delivery/http_test"
)

func Run() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server
	err := http2.StartServer(&cnf.Server)
	if err != nil {
		return err
	}

	return nil
}

func RunServerTest() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server for testing
	err := http_test2.StartServer(&cnf.ServerTest)
	if err != nil {
		return err
	}

	return nil
}
