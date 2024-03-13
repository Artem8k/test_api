package service

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db *mongo.Client
}

func New(db *mongo.Client) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetJwtPair(w http.ResponseWriter, req *http.Request) http.Request {
	fmt.Println(req)
	return *req
}

func (s *Service) UpdateAccessToken(w http.ResponseWriter, req *http.Request) http.Request {
	fmt.Println(req)
	return *req
}
