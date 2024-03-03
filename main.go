package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cache *Cache

func main() {
	// подключение к серверу NATS Streaming
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// подключение к бд
	connStr := "user=postgres password=2324 dbname=order sslmode=disable"
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	cache = NewCache()

	// подписка на канал
	subscribe(nc, err, db)

	// достаем данные из бд
	rows, err := db.Query(context.Background(), "SELECT * FROM order_info")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// записываем данные в кеш
	writeCache(rows)

	fmt.Println("Listening for messages...")

	// старт http сервера
	StartServer(cache)

	// ожидание завершения программы
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down...")
}

func subscribe(nc *nats.Conn, err error, db *pgx.Conn) {
	_, err = nc.Subscribe("purchases", func(m *nats.Msg) {
		// принимает в структуру переданный json по байтам
		var orderData OrderFields
		if err := json.Unmarshal(m.Data, &orderData); err != nil {
			log.Printf("Failed to read json: %v\n", err)
			return
		}
		// данные отправляются в кеш
		cache.Set(orderData.OrderUid, orderData)
		// запись в бд
		_, err = db.Exec(context.Background(), "INSERT INTO order_info (order_uid, info) VALUES ($1, $2)",
			orderData.OrderUid, orderData)
		if err != nil {
			log.Printf(err.Error())
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
}

func writeCache(rows pgx.Rows) {
	for rows.Next() {
		var order OrderFields
		var uid string
		var oStr string
		err := rows.Scan(&uid, &oStr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = json.Unmarshal([]byte(oStr), &order)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cache.Set(order.OrderUid, order)
	}
}
