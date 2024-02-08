package handlers

import (
	"WB-L0-n/my_cache"
	"github.com/gin-gonic/gin"
	"log"
)

func GetOrder(c *gin.Context) {
	uid := c.Query("order-uid")

	intrf, _ := my_cache.OrderCache.Get(uid)
	log.Println("!!!!!!!!!!!!", intrf)

	c.JSON(200, intrf)
}
