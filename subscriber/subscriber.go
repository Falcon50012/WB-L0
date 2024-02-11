package subscriber

import (
	"WB-L0-n/database"
	"WB-L0-n/model"
	"WB-L0-n/my_cache"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"log"
)

func InMemory(order *model.Order) {
	my_cache.OrderCache.Set(order.OrderUID, order, cache.DefaultExpiration)
	fmt.Println("Cache Content:")
	for key, value := range my_cache.OrderCache.Items() {
		fmt.Printf("Key: %s, Value: %+v\n", key, value.Object)
	}
}

var ctx, cancel = context.WithCancel(context.Background())

func Subscribe() error {
	clusterID := "test-cluster"
	clientID := "WB-sub"
	URL := "nats://localhost:4223"
	subj := "OrderData"

	nc, err := nats.Connect(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(subj, func(msg *stan.Msg) {
		var newOrder model.Order
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Println(err)
			return
		}
		database.SaveOrderToDB(&newOrder)
		InMemory(&newOrder)
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	fmt.Println("Subscribed to NATS Streaming. Waiting for messages...")

	<-ctx.Done()

	return nil
}

func StopSubscribe() {
	cancel()
}
