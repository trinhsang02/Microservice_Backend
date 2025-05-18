package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"common-service/api"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/yourusername/yourrepo/db/sqlc"
)

func checkDatabaseConnection(pool *pgxpool.Pool) error {
	// Try to ping the database
	err := pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Get connection stats
	stats := pool.Stat()
	log.Printf("Database connection stats:")
	log.Printf("- Total connections: %d", stats.TotalConns())
	log.Printf("- Acquired connections: %d", stats.AcquiredConns())
	log.Printf("- Max connections: %d", stats.MaxConns())

	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	log.Printf("Connecting to database with DSN: %s", dsn)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Check database connection
	if err := checkDatabaseConnection(pool); err != nil {
		log.Fatalf("Database connection check failed: %v", err)
	}
	log.Println("Successfully connected to database!")

	// Initialize queries and repository
	repo := sqlc.NewRepository(sqlc.New(pool))

	// Initialize API handler
	handler := api.NewHandler(repo)

	// Setup router
	router := api.SetupRouter(handler)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
