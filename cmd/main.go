package main

import (
	"WB-L0-n/database"
	"WB-L0-n/my_cache"
	"WB-L0-n/router"
	"WB-L0-n/subscriber"
	"log"
)

func init() {
	my_cache.InitCache()
	database.InitDBConnectionPool()
	err := database.CreateTables()
	if err != nil {
		log.Fatal(err)
	}
	subscriber.RestoreCache()
}

func main() {
	log.Println(my_cache.OrderCache)
	r := router.InitRoutes()

	log.Println("Running Server and Subscriber")

	go func() {
		subscriber.Subscribe()
	}()

	//subscriber.Subscribe()

	if err := r.Run(); err != nil {
		log.Fatalf("Ошибка при запуске веб-сервера: %v", err)
	}
	subscriber.StopSubscribe()
}

//err := r.Run("127.0.0.1:8080")
//if err != nil {
//	return
//}
//}
