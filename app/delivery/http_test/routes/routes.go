package routes

import (
	"broken-link-checker/app/delivery/http_test/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRoutes() *gin.Engine {
	gin.SetMode(viper.GetString("server.mode"))

	r := gin.New()
	r.Use(
		// gin.Logger(),
		gin.Recovery(),
	)

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
	redirectGroup := r.Group("/redirect")
	{
		redirectGroup.GET("*any", response.Redirect)
	}

	r.NoRoute(func(c *gin.Context) {
		response.Success(c)
	})

	return r
}
