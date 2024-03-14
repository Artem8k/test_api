package service

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"rest-api/src/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshToken struct {
	exp time.Time
}

type Service struct {
	db *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) AddDataToDB(w http.ResponseWriter, req *http.Request) {
	hash, err := s.generateRefreshToken()

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = s.db.Client.Database("local").Collection("refreshToken").InsertOne(s.db.Ctx, bson.D{{Key: "hash", Value: hash}, {Key: "guid", Value: req.URL.Query().Get("guid")}})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Можно предположить что запрос на этот метод приходит после успешной валидации данных пользователя другим сервисом
func (s *Service) GetJwtPair(w http.ResponseWriter, req *http.Request) (*TokenPair, error) {
	// TODO: изменить аргументы функции, чтобы не передавать w и req
	// здесь получать сразу guid из req
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"guid": req.URL.Query().Get("guid"),
			"exp":  time.Now().Add(15 * time.Minute).Unix(),
		})

	signedAccessToken, err := accessToken.SignedString([]byte("TOP_SECRET")) // тут в идеале нужно добавить получение ключа из переменной окружения
	if err != nil {
		return nil, err
	}

	hash, err := s.generateRefreshToken()

	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "guid", Value: req.URL.Query().Get("guid")}}
	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "hash", Value: hash}, {Key: "guid", Value: req.URL.Query().Get("guid")}},
	}}
	_, err = s.db.Client.Database("local").Collection("refreshToken").UpdateOne(s.db.Ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: base64.RawStdEncoding.EncodeToString(hash),
	}, nil
}

func (s *Service) UpdateAccessToken(w http.ResponseWriter, req *http.Request) {
	// TODO: из req headers auth получить access токен
	// провалидировать refresh token из body: посмотреть exp
	// достать из access токена guid
	// найти по guid запись в бд
	// провалидировать refresh token из body: сравнить с refresh в бд
	// создать оба токена
	// обновить refresh в бд
	fmt.Println(req)
}

func (s *Service) generateRefreshToken() ([]byte, error) {
	var sec int64 = 15552000 //180 days
	var token = refreshToken{
		exp: time.Unix(sec, 0),
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprint(token)), 10)

	if err != nil {
		return nil, err
	}

	return hash, nil
}
