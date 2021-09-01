package app

import (
	"broken-link-checker/app/config"
	"broken-link-checker/app/delivery/http_test"
	"broken-link-checker/app/service/linkChecker"
	"fmt"
	"time"
)

func Run() error {
	// Get the server settings
	cnf := config.Get()

	// Start the server for testing
	err := http_test.Start(&cnf.ServerTest)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	start := time.Now()
	checker := linkChecker.Checker{}
	err = checker.Run("http://localhost:8080/", 3)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	fmt.Println("Time spent: ", duration)

	return nil
}
