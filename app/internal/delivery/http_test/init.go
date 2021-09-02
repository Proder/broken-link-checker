package http_test

import (
	routes2 "broken-link-checker/app/internal/delivery/http_test/routes"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
	Mode string
}

func StartServer(cnf *Config) error {
	gin.SetMode(cnf.Mode)

	// Declaring routes
	rts := routes2.InitRoutes()

	// Starting the server
	if err := rts.Run(cnf.Host + ":" + cnf.Port); err != nil {
		return err
	}

	return nil
}
