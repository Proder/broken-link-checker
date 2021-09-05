package routes

import (
	"broken-link-checker/internal/delivery/http/api/v1/checker"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(
		// gin.Logger(),
		gin.Recovery(),
	)

	r.Use(static.Serve("/", static.LocalFile("./web", true)))

	apiGroup := r.Group("/api")
	{
		v1 := apiGroup.Group("/v1")
		{
			v1.POST("search-broken-links", checker.SearchBrokenLinks)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, nil)
	})

	return r
}
