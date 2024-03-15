package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api/src/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type refreshDocument struct {
	Hash string `json:"hash"`
	Guid string `json:"guid"`
	Exp  int64  `json:"exp"`
}

type httpError struct {
	Error  string `json:"error"`
	Status int    `json:"statusCode"`
}

type AccessTokenClaims struct {
	Guid string `json:"guid"`
	Exp  int64  `json:"exp"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Service struct {
	db *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetJwtPair(w http.ResponseWriter, guid string) *TokenPair {
	if guid == "" {
		s.errorHandler(w, httpError{
			Error:  "Bad Request",
			Status: 400,
		})
	}

	// создаю jwt токен c payload в котором есть guid и expiration
	signedAccessToken, err := s.generateJwt(guid)
	if err != nil {
		e := httpError{Error: err.Error(), Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	// создаю хеш из guid
	hash, err := s.generateRefreshToken(guid)

	if err != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	exp := time.Unix(time.Now().Unix()+15552000, 0).Unix() //+180 days
	filter := bson.D{{Key: "guid", Value: guid}}
	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "hash", Value: string(hash)}, {Key: "guid", Value: guid}, {Key: "exp", Value: exp}},
	}}
	_, err = s.db.Client.Database("local").Collection("refreshToken").UpdateOne(s.db.Ctx, filter, update)

	if err != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: base64.RawStdEncoding.EncodeToString(hash),
	}
}

func (s *Service) UpdateAccessToken(w http.ResponseWriter, accessToken string, refreshToken string) *TokenPair {
	// расшифровываю токен
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("TOPSECRET"), nil
	})
	if err != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	// достаю payload
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	// ищу по guid из токена запись в таблице и декодирую запись
	var r refreshDocument
	document := s.db.Client.Database("local").Collection("refreshToken").FindOne(s.db.Ctx, bson.D{{Key: "guid", Value: claims.Guid}})
	document.Decode(&r)
	if document.Err() != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(document.Err())
		return nil
	}

	// сравниваю hash из body с hash в бд
	refreshHash, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		e := httpError{Error: "Invalid refresh token", Status: 404}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}
	if !bytes.Equal(refreshHash, []byte(r.Hash)) {
		e := httpError{Error: "Invalid refresh token", Status: 401}
		s.errorHandler(w, e)
		return nil
	}

	// проверяю не истек ли токен
	t := time.Now().Unix()
	exp := r.Exp
	if t > exp {
		e := httpError{Error: "Refresh token expired", Status: 401}
		s.errorHandler(w, e)
		return nil
	}

	signedAccessToken, err := s.generateJwt(claims.Guid)
	if err != nil {
		e := httpError{Error: err.Error(), Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	hash, err := s.generateRefreshToken(claims.Guid)

	if err != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	expiration := time.Unix(time.Now().Unix()+15552000, 0).Unix() //+180 days
	filter := bson.D{{Key: "guid", Value: claims.Guid}}
	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "hash", Value: string(hash)}, {Key: "guid", Value: claims.Guid}, {Key: "exp", Value: expiration}},
	}}
	_, err = s.db.Client.Database("local").Collection("refreshToken").UpdateOne(s.db.Ctx, filter, update)

	if err != nil {
		e := httpError{Error: "Internal server error", Status: 500}
		s.errorHandler(w, e)
		fmt.Println(err)
		return nil
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: base64.RawStdEncoding.EncodeToString(hash),
	}
}

func (s *Service) generateRefreshToken(guid string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(guid), 5)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (s *Service) generateJwt(guid string) (string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512,
		AccessTokenClaims{
			Guid: guid,
			Exp:  time.Now().Add(15 * time.Minute).Unix(),
		})

	signedAccessToken, err := accessToken.SignedString([]byte("TOPSECRET"))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func (s *Service) errorHandler(w http.ResponseWriter, e httpError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
