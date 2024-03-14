package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/src/service"

	"github.com/gorilla/mux"
)

// ИНИЦИАЛИЗИРОВАТЬ БД
// СДЕЛАТЬ APP ФАЙЛ ДЛЯ ЗАПУСКА БД СЕРВИСА И ХЕНДЛЕРОВ
// ЗАПУСТИТЬ ЕГО В MAIN.GO ФАЙЛЕ

type Handlers struct {
	service *service.Service
}

func New(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) InitHandlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", h.GetJwtPair).Methods("GET")
	router.HandleFunc("/refresh", h.UpdateAccessToken).Methods("POST")
	router.HandleFunc("/testData", h.service.AddDataToDB).Methods("POST")

	return router
}

func (h *Handlers) GetJwtPair(w http.ResponseWriter, req *http.Request) {
	res, _ := h.service.GetJwtPair(w, req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handlers) UpdateAccessToken(w http.ResponseWriter, req *http.Request) {
	h.service.UpdateAccessToken(w, req)
}

func (h *Handlers) AddDataToDB(w http.ResponseWriter, req *http.Request) {
	h.service.AddDataToDB(w, req)
	//w.Write(res)
}
