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
	// Lấy các thông tin kết nối từ biến môi trường
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "microservices")

	// Tạo connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Kết nối đến database
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("không thể kết nối đến database: %v", err)
	}

	// Kiểm tra kết nối
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("không thể ping đến database: %v", err)
	}

	log.Println("Kết nối PostgreSQL thành công - relayer-service")
	return nil
}

// CloseDB đóng kết nối với database
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Đã đóng kết nối PostgreSQL - relayer-service")
	}
}

// getEnv lấy giá trị biến môi trường, nếu không có trả về giá trị mặc định
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
