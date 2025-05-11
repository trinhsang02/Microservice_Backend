package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB khởi tạo kết nối đến PostgreSQL
func InitDB() error {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "microservices")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect database
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("can not connect database: %v", err)
	}

	// Kiểm tra kết nối
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("can not ping to database: %v", err)
	}

	log.Println("Connect PostgreSQL successful - blockchain-listener")
	return nil
}

// CloseDB 
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Đã đóng kết nối PostgreSQL - blockchain-listener")
	}
}

// getEnv 
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
