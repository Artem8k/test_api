package handlers

import (
	"encoding/json"
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
	router.HandleFunc("/login", h.GetJwtPair).Methods("GET")
	router.HandleFunc("/refresh", h.UpdateAccessToken).Methods("POST")

	return router
}

func (h *Handlers) GetJwtPair(w http.ResponseWriter, req *http.Request) {
	guid := req.URL.Query().Get("guid")
	res := h.service.GetJwtPair(w, guid)

	if res != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else {
		return
	}
}

func (h *Handlers) UpdateAccessToken(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	accessToken := strings.Split(authHeader, "Bearer ")[1]

	var token RefreshReq
	json.NewDecoder(req.Body).Decode(&token)
	res := h.service.UpdateAccessToken(w, accessToken, token.RefreshToken)

	if res != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	} else {
		return
	}
}
