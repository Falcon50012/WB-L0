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

//var Item sub_models.ItemInfo

func InMemory(order *model.Order) {
	my_cache.OrderCache.Set(order.OrderUID, order, cache.DefaultExpiration)
	//fmt.Println("Cache Content:")
	//for key, value := range my_cache.OrderCache.Items() {
	//	fmt.Printf("Key: %s, Value: %+v\n", key, value.Object)
	//}
}

func RestoreCache() {

	rows, err := database.DBPool.Query(context.Background(), "SELECT orders.order_uid, orders.track_number, orders.entry, orders.locale, orders.internal_signature, orders.customer_id, orders.delivery_service, orders.shardkey, orders.sm_id, orders.date_created, orders.oof_shard, deliveries.name, deliveries.phone, deliveries.zip, deliveries.city, deliveries.address, deliveries.region, deliveries.email, payments.transaction, payments.request_id, payments.currency, payments.provider, payments.amount, payments.payment_dt, payments.bank, payments.delivery_cost, payments.goods_total, payments.custom_fee, items.chrt_id, items.track_number, items.price, items.rid, items.name, items.sale, items.size, items.total_price, items.nm_id, items.brand, items.status FROM orders LEFT JOIN public.deliveries ON orders.order_uid = public.deliveries.order_uid LEFT JOIN public.payments ON orders.order_uid = public.payments.order_uid LEFT JOIN public.items ON orders.order_uid = public.items.order_uid LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	my_cache.CacheMx.Lock()
	defer my_cache.CacheMx.Unlock()

	for rows.Next() {
		var order model.Order
		var item model.ItemInfo

		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
			&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)

		if err != nil {
			log.Println(err)
			continue
		}

		co, found := my_cache.OrderCache.Get(order.OrderUID)
		if found {
			existingOrder := co.(model.Order)
			existingOrder.Items = append(existingOrder.Items, item)
			my_cache.OrderCache.Set(order.OrderUID, existingOrder, cache.DefaultExpiration)
		} else {
			order.Items = []model.ItemInfo{item}
			my_cache.OrderCache.Set(order.OrderUID, order, cache.DefaultExpiration)
		}
	}
	fmt.Println("Cache Content:")
	for key, value := range my_cache.OrderCache.Items() {
		fmt.Printf("Key: %s, Value: %+v\n", key, value.Object)
	}
}

var ctx, cancel = context.WithCancel(context.Background())

func Subscribe() {
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
}

func StopSubscribe() {
	cancel()
}

//func Subscribe() {
//
//	var (
//		clusterID   string
//		clientID    string
//		URL         string
//		subj        string
//		qgroup      string
//		durable     string
//		unsubscribe bool
//	)
//
//	clusterID = "test-cluster"
//	clientID = "WB-pub"
//	URL = "nats://localhost:4223"
//	subj = "OrderData"
//
//	nc, err := nats.Connect(URL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer nc.Close()
//
//	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
//		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
//			log.Fatalf("Connection lost, reason: %v", reason)
//		}))
//	if err != nil {
//		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
//	}
//	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clientID)
//	defer sc.Close()
//
//	mcb := func(msg *stan.Msg) {
//		var newOrder model.Order
//
//		data := msg.Data
//
//		err = json.Unmarshal(data, &newOrder)
//		if err != nil {
//			log.Println(err)
//		}
//		//fmt.Printf("%+v\n", newOrder)
//		database.SaveOrderToDB(&newOrder)
//
//		InMemory(&newOrder)
//	}
//
//	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, stan.DurableName(durable))
//	if err != nil {
//		sc.Close()
//		log.Fatal(err)
//	}
//
//	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj, clientID, qgroup, durable)
//
//	signalChan := make(chan os.Signal, 1)
//	cleanupDone := make(chan bool)
//	signal.Notify(signalChan, os.Interrupt)
//	go func() {
//		for range signalChan {
//			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
//			// Do not unsubscribe a durable on exit, except if asked to.
//			if durable == "" || unsubscribe {
//				sub.Unsubscribe()
//			}
//			sc.Close()
//			cleanupDone <- true
//		}
//	}()
//	<-cleanupDone
//}