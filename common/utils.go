package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type Envelop = map[string]any

type DatabaseConfig struct {
	dbUsername string
	dbPassword string
	dbHost     string
	dbPort     int
	dbName     string
}

func GetDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load("../common/.env")
	if err != nil {
		return nil, err
	}
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		dbPort = 5432
	}
	return &DatabaseConfig{
		dbUsername: dbUsername,
		dbPassword: dbPassword,
		dbHost:     dbHost,
		dbName:     dbName,
		dbPort:     dbPort,
	}, nil
}

func ConnectToDB() (*sql.DB, error) {
	dbConfig, err := GetDatabaseConfig()
	if err != nil {
		return nil, err
	}
	connectionStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		dbConfig.dbUsername,
		dbConfig.dbPassword,
		dbConfig.dbHost,
		dbConfig.dbPort,
		dbConfig.dbName,
	)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(Envelop{"error": "Bad JSON format"})
		w.Write(errMsg)
	}
	w.Write(jsonData)
}

func ReadJSON(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}
