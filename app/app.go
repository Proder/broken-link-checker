package app

import (
	"broken-link-checker/app/config"
	"broken-link-checker/app/delivery/http"
	"broken-link-checker/app/delivery/http_test"
)

func Run() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server
	err := http.StartServer(&cnf.Server)
	if err != nil {
		return err
	}

	return nil
}

func RunServerTest() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server for testing
	err := http_test.StartServer(&cnf.ServerTest)
	if err != nil {
		return err
	}

	return nil
}
