package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api/src/database"
)

type httpError struct {
	Error  string `json:"error"`
	Status int    `json:"statusCode"`
}

type Order struct {
	Product_Id   int64  `json:"product_id" db:"product_id"`
	Product_Name string `json:"product_name" db:"product_name"`
	Shelf_Id     string `json:"shelf_id" db:"shelf_id"`
	Shelf_Name   string `json:"shelf_name" db:"shelf_name"`
	Quantity     int64  `json:"quantity" db:"quantity"`
	Order_Id     int64  `json:"order_id" db:"order_id"`
}

type Service struct {
	db *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetOrders(w http.ResponseWriter, orders []string) []*Order {
	var o []*Order
	err := s.db.Client.Select(&o,
		`SELECT p.id as product_id, s.id as shelf_id, s.name as shelf_name, 
		p.name as product_name, p_o.quantity as quantity, 
		p_o.order_id as order_id FROM product AS p
		
		JOIN product_in_order AS p_o ON p.id = p_o.product_id
		JOIN shelf as s ON s.id = p_o.shelf_id
		
		WHERE p_o.order_id = any ($1)
		
		GROUP BY p.id, s.id, s.name, p_o.quantity, p.name, order_id
		ORDER BY shelf_name, p_o.order_id`, orders)

	if err != nil {
		fmt.Println(err)
		s.errorHandler(w, httpError{
			Status: 500,
			Error:  err.Error(),
		})
		return nil
	}
	return o
}

func (s *Service) errorHandler(w http.ResponseWriter, e httpError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
