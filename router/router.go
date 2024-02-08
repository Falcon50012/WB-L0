package router

import (
	"WB-L0-n/handlers"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"path"
)

func InitRoutes() *gin.Engine {
	frontendDir := "front"
	router := gin.New()
	router.Use(static.Serve("/", static.LocalFile(frontendDir, true)))

	router.GET("get-order", handlers.GetOrder)

	router.NoRoute(func(c *gin.Context) {
		c.File(path.Join(frontendDir, "index.html"))
	})

	return router
}
