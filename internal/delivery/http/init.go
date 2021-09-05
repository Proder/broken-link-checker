package http

import (
	"broken-link-checker/internal/delivery/http/routes"
	"fmt"

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
	rts := routes.InitRoutes()

	// Starting the server
	if err := rts.Run(cnf.Host + ":" + cnf.Port); err != nil {
		return fmt.Errorf("rts.Run failed: %w", err)
	}

	return nil
}
