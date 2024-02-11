package database

import (
	"WB-L0-n/model"
	"WB-L0-n/my_cache"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/patrickmn/go-cache"
	"log"
)

func SaveOrderToDB(order *model.Order) {

	_, err := DBPool.Exec(context.Background(),
		"INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		log.Fatalf("ERROR into orders %v", err)
	}

	_, err = DBPool.Exec(context.Background(),
		"INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Fatalf("ERROR into deliveries %v", err)
	}

	_, err = DBPool.Exec(context.Background(),
		"INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Fatalf("ERROR into payments %v", err)
	}

	for _, Item := range order.Items {
		_, err = DBPool.Exec(context.Background(),
			"INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUID, Item.ChrtID, Item.TrackNumber, Item.Price, Item.Rid, Item.Name, Item.Sale, Item.Size,
			Item.TotalPrice, Item.NmID, Item.Brand, Item.Status)
		if err != nil {
			log.Fatalf("ERROR into items %v", err)
		}
	}
}

func RestoreCache() error {
	rows, err := DBPool.Query(context.Background(), "SELECT orders.order_uid, orders.track_number, orders.entry, orders.locale, orders.internal_signature, orders.customer_id, orders.delivery_service, orders.shardkey, orders.sm_id, orders.date_created, orders.oof_shard, deliveries.name, deliveries.phone, deliveries.zip, deliveries.city, deliveries.address, deliveries.region, deliveries.email, payments.transaction, payments.request_id, payments.currency, payments.provider, payments.amount, payments.payment_dt, payments.bank, payments.delivery_cost, payments.goods_total, payments.custom_fee, items.chrt_id, items.track_number, items.price, items.rid, items.name, items.sale, items.size, items.total_price, items.nm_id, items.brand, items.status FROM orders LEFT JOIN public.deliveries ON orders.order_uid = public.deliveries.order_uid LEFT JOIN public.payments ON orders.order_uid = public.payments.order_uid LEFT JOIN public.items ON orders.order_uid = public.items.order_uid ORDER BY date_created DESC LIMIT 12")
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

	return nil
}

func GetOrderByID(orderUID string) (model.Order, error) {
	var order model.Order

	row := DBPool.QueryRow(context.Background(), "SELECT orders.order_uid, orders.track_number, orders.entry, orders.locale, orders.internal_signature, orders.customer_id, orders.delivery_service, orders.shardkey, orders.sm_id, orders.date_created, orders.oof_shard, deliveries.name, deliveries.phone, deliveries.zip, deliveries.city, deliveries.address, deliveries.region, deliveries.email, payments.transaction, payments.request_id, payments.currency, payments.provider, payments.amount, payments.payment_dt, payments.bank, payments.delivery_cost, payments.goods_total, payments.custom_fee FROM orders LEFT JOIN public.deliveries ON orders.order_uid = public.deliveries.order_uid LEFT JOIN public.payments ON orders.order_uid = public.payments.order_uid WHERE orders.order_uid = $1", orderUID)

	my_cache.CacheMx.Lock()
	defer my_cache.CacheMx.Unlock()

	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
		&order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
		&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return order, err
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Unable to find order")
			return order, nil
		}
		log.Printf("Error executing request: %v", err)
		return order, err
	}

	rows, err := DBPool.Query(context.Background(), "SELECT items.chrt_id, items.track_number, items.price, items.rid, items.name, items.sale, items.size, items.total_price, items.nm_id, items.brand, items.status FROM items WHERE items.order_uid = $1", orderUID)
	if err != nil {
		log.Printf("Error when executing a request to receive items: %v", err)
		return order, err
	}
	defer rows.Close()

	var items []model.ItemInfo
	for rows.Next() {
		var item model.ItemInfo
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Printf("Error while scanning item: %v", err)
			return order, err
		}
		items = append(items, item)
	}

	order.Items = items

	return order, nil
}
