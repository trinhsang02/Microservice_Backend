package api

import (
	"common-service/database"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// SetupRouter configures and returns the API router
func SetupRouter(repo *database.Repository) *mux.Router {
	handler := NewHandler(repo)
	router := mux.NewRouter()

	// Health check route
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// KYC routes
	apiRouter.HandleFunc("/kyc/submit", handler.SubmitKYC).Methods("POST")
	apiRouter.HandleFunc("/kyc/citizen/{citizen_id}", handler.GetKYCByCitizenID).Methods("GET")
	apiRouter.HandleFunc("/kyc/wallet/{wallet_address}", handler.GetKYCByWalletAddress).Methods("GET")

	// NFT routes
	apiRouter.HandleFunc("/kyc/mint-nft", handler.MintKYCNFT).Methods("POST")

	// Log all requests
	router.Use(loggingMiddleware)

	return router
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
