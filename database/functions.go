package database

import (
	"WB-L0-n/model"
	"context"
	"log"
)

func SaveOrderToDB(order *model.Order) {
	//log.Printf("Inserting order with values: %+v\n", order)

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
