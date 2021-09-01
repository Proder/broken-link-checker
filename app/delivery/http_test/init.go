package http_test

import (
	"broken-link-checker/app/delivery/http_test/routes"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
}

func Start(cnf *Config) error {
	gin.SetMode("release")

	// Declaring routes
	rts := routes.InitRoutes()

	// Starting the server
	go func() {
		_ = rts.Run(cnf.Host + ":" + cnf.Port)
	}()

	return nil
}
