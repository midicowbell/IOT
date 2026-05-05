package main

import (
	"fmt"
	"iot/handlers"
	"iot/server"
	"iot/service"
	"iot/storage"
)

func main() {
	repo := storage.NewMemStorage()
	stockService := service.NewStockService(repo)
	httpHandlers := handlers.NewHTTPHandlers(stockService)
	srv := server.NewServer(httpHandlers)
	if err := srv.StartServer(); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
