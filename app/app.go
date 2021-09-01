package app

import (
	"broken-link-checker/app/config"
	"broken-link-checker/app/delivery/http_test"
	"broken-link-checker/app/service/linkChecker"
)

func Run() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server for testing
	err := http_test.Start(&cnf.ServerTest)
	if err != nil {
		return err
	}

	err = linkChecker.Run("http://localhost:8080/", 3)
	if err != nil {
		return err
	}
	return nil
}
