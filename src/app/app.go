package app

import (
	"net/http"
	"rest-api/src/database"
	handlers "rest-api/src/handelrs"
	"rest-api/src/service"
)

func MustRun() {
	db := database.MustRun()

	service := service.New(db)

	handlers := handlers.New(service)

	http.ListenAndServe(":3000", handlers.InitHandlers())
}
