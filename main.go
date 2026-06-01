package main

import (
	"context"
	"fmt"
	"iot/handlers"
	"iot/server"
	"iot/service"
	"iot/storage"
)

func main() {
	ctx := context.Background()
	url := "postgres://postgres:12345@localhost:5432/warehouse_iot"
	repo, err := storage.NewPostgresStorage(ctx, url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer repo.Close()
	stockService := service.NewStockService(repo)
	httpHandlers := handlers.NewHTTPHandlers(stockService)
	srv := server.NewServer(httpHandlers)
	if err := srv.StartServer(); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
