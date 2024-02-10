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
		log.Fatalf("Error creating tables: %v", err)
	}

	err = subscriber.RestoreCache()
	if err != nil {
		log.Printf("Error while restoring cache: %v", err)
	}
}

func main() {
	log.Println("Running Server and Subscriber")

	go func() {
		err := subscriber.Subscribe()
		if err != nil {
			log.Printf("Error when starting subscriber: %v", err)
		}
	}()

	r := router.InitRoutes()

	err := r.Run()
	if err != nil {
		log.Fatalf("Error starting web-server: %v", err)
	}

	subscriber.StopSubscribe()
}
