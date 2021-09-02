package routes

import (
	"broken-link-checker/app/internal/delivery/http_test/response"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(
		// gin.Logger(),
		gin.Recovery(),
	)

	r.GET("/", response.Success)

	successGroup := r.Group("/success")
	{
		successGroup.GET("*any", response.Success)
	}

	errGroup := r.Group("/error")
	{
		urlGroup := errGroup.Group("/url")
		{
			urlGroup.GET("*any", response.ErrorUrl)
		}
		serverGroup := errGroup.Group("/server")
		{
			serverGroup.GET("*any", response.ErrorServer)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		response.ErrorUrl(c)
	})

	return r
}
