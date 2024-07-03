package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ContextKey string

const ContextUserIdKey ContextKey = "userId"

type Envelop = map[string]any

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		errMsg, _ := json.Marshal(Envelop{"error": "bad json format"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg)
		return
	}
	w.WriteHeader(status)
	w.Write(jsonData)
}

func ReadJSON(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func GetEnvValByKey(key string, fallback string) string {
	godotenv.Load(".env")
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ConnectToDatabase(ctx context.Context, driverName string, connectionStr string) (*sql.DB, error) {
	db, err := sql.Open(driverName, connectionStr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GenerateNanoid() string {
	nanoid, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	return nanoid
}

func GenerateSlugName(name string) string {
	slug := strings.ReplaceAll(name, " ", "-")
	slug = strings.ToLower(slug)
	return slug
}

func SignToken(data map[string]any, expiredInMinutes time.Duration, secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().Add(expiredInMinutes * time.Minute).UTC().Unix()
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	invalidClaim := token.Claims.Valid()
	if invalidClaim != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
