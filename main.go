package main

import (
	"fmt"
	"net/http"
	"rest-api/src/database"
	"rest-api/src/handlers"
	"rest-api/src/service"
)

func main() {
	fmt.Println("Starting an application")

	db := database.MustRun()

	service := service.New(db)

	handlers := handlers.New(service)

	err := http.ListenAndServe(":3000", handlers.InitHandlers())

	if err != nil {
		panic(err)
	}
	fmt.Println("Application stopped")
}
