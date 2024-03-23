package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api/src/service"
	"strings"

	"github.com/gorilla/mux"
)

type RefreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

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
	router.HandleFunc("/getOrders", h.GetOrders).Methods("GET")
	return router
}

func (h *Handlers) GetOrders(w http.ResponseWriter, req *http.Request) {
	ordersString := req.URL.Query().Get("orders")
	splitFunc := func(r rune) bool {
		return strings.ContainsRune("[],", r)
	}
	orders := strings.FieldsFunc(ordersString, splitFunc)
	res := h.service.GetOrders(w, orders)

	if res != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(res)
		json.NewEncoder(w).Encode(res)
	} else {
		return
	}
}
