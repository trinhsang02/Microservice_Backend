package main

import (
	"common-service/api"
	"common-service/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	// Create repository and setup router
	repo := database.NewRepository()
	router := api.SetupRouter(repo)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure server
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	log.Printf("API Service starting on port %s", port)
	log.Fatal(server.ListenAndServe())
}
