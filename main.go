package main

import (
	"context"
	"fmt"
	"iot/handlers"
	"iot/server"
	"iot/service"
	"iot/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	url := "postgres://postgres:12345@localhost:5432/warehouse_iot"
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Println(err)
	}
	repo := storage.NewDataStorage(pool)
	stockService := service.NewStockService(repo)
	httpHandlers := handlers.NewHTTPHandlers(stockService)
	srv := server.NewServer(httpHandlers)
	if err := srv.StartServer(); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
