package handlers

import (
	"WB-L0-n/database"
	"WB-L0-n/my_cache"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"log"
)

func GetOrder(c *gin.Context) {
	uid := c.Query("order-uid")
	if uid == "" {
		c.JSON(400, gin.H{"error": "order_uid not found or empty"})
		return
	}

	order, found := my_cache.OrderCache.Get(uid)
	if found {
		log.Println("Order is found in cache", order)
		c.JSON(200, order)
		return
	}

	log.Println("Order is not found in cache, retrieving from database...")
	dbOrder, err := database.GetOrderByID(uid)
	if err != nil {
		log.Printf("Error while fetching order from database: %v", err)
		c.JSON(500, gin.H{"error": "failed to fetch order from database"})
		return
	}

	my_cache.OrderCache.Set(uid, dbOrder, cache.DefaultExpiration)
	log.Println("Order is retrieved from database and added to cache", dbOrder)

	c.JSON(200, dbOrder)
}
