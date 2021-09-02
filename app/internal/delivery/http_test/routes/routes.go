package routes

import (
	response2 "broken-link-checker/app/internal/delivery/http_test/response"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(
		// gin.Logger(),
		gin.Recovery(),
	)

	r.GET("/", response2.Success)

	successGroup := r.Group("/success")
	{
		successGroup.GET("*any", response2.Success)
	}

	errGroup := r.Group("/error")
	{
		urlGroup := errGroup.Group("/url")
		{
			urlGroup.GET("*any", response2.ErrorUrl)
		}
		serverGroup := errGroup.Group("/server")
		{
			serverGroup.GET("*any", response2.ErrorServer)
		}
	}
	redirectGroup := r.Group("/redirect")
	{
		redirectGroup.GET("*any", response2.Redirect)
	}

	r.NoRoute(func(c *gin.Context) {
		response2.ErrorUrl(c)
	})

	return r
}
