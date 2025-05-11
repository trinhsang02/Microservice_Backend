package api

import (
	"common-service/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Handler struct holds dependencies for API handlers
type Handler struct {
	repo *database.Repository
}

// NewHandler creates a new Handler instance
func NewHandler(repo *database.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// HealthCheck provides a simple health check endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, database.ApiResponse{
		Success: true,
		Message: "API service is running",
	})
}

// KYCRequest represents the request payload for KYC submission
type KYCRequest struct {
	CitizenID       string `json:"citizen_id"`
	WalletAddress   string `json:"wallet_address"`
	FullName        string `json:"full_name"`
	PhoneNumber     string `json:"phone_number"`
	DateOfBirth     string `json:"date_of_birth"` // Format: YYYY-MM-DD
	Nationality     string `json:"nationality"`
	KYCVerifiedAt   string `json:"kyc_verified_at,omitempty"` // Format: YYYY-MM-DD HH:MM:SS or empty
	Verifier        string `json:"verifier,omitempty"`
	WalletSignature string `json:"wallet_signature,omitempty"`
}

// SubmitKYC handles submission of KYC information with new schema
func (h *Handler) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	var request KYCRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Log the received request for debugging
	requestJSON, _ := json.Marshal(request)
	log.Printf("Received KYC submission: %s", string(requestJSON))

	// Validate required fields
	if request.CitizenID == "" || request.WalletAddress == "" ||
		request.FullName == "" || request.PhoneNumber == "" ||
		request.DateOfBirth == "" || request.Nationality == "" || 
		request.WalletSignature == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields for KYC submission")
		return
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format for date of birth. Use YYYY-MM-DD format")
		return
	}

	// Create KYC information object
	kycInfo := database.KYC{
		CitizenID:       request.CitizenID,
		WalletAddress:   request.WalletAddress,
		FullName:        request.FullName,
		PhoneNumber:     request.PhoneNumber,
		DateOfBirth:     dob,
		Nationality:     request.Nationality,
		IsActive:        true,
		WalletSignature: request.WalletSignature,
	}

	// Always set verification time to current time
	verifiedTime := time.Now()
	kycInfo.KYCVerifiedAt = &verifiedTime

	// Always set verifier ( default to "system")
	kycInfo.Verifier = request.Verifier
	if kycInfo.Verifier == "" {
		kycInfo.Verifier = "system"
	}

	// Log the KYC object being stored
	log.Printf("Storing KYC: CitizenID=%s, VerifiedAt=%v, Verifier=%s",
		kycInfo.CitizenID,
		kycInfo.KYCVerifiedAt,
		kycInfo.Verifier)

	// Submit KYC information
	err = h.repo.SubmitKYC(kycInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to submit KYC information: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, database.ApiResponse{
		Success: true,
		Message: "KYC information submitted successfully",
		Data:    kycInfo,
	})
}

// MintKYCNFT handles minting an NFT for a user with verified KYC
func (h *Handler) MintKYCNFT(w http.ResponseWriter, r *http.Request) {
	var request struct {
		WalletAddress string `json:"wallet_address"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if request.WalletAddress == "" {
		respondWithError(w, http.StatusBadRequest, "Wallet address is required")
		return
	}

	// Get KYC information for the wallet address
	kycInfo, err := h.repo.GetKYCByWalletAddress(request.WalletAddress)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "KYC information not found for this wallet address")
		return
	}

	// Verify that KYC is active
	if !kycInfo.IsActive {
		respondWithError(w, http.StatusBadRequest, "KYC is not active for this wallet address")
		return
	}

	// Mint NFT
	err = h.repo.MintKYCNFT(*kycInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to mint NFT: "+err.Error())
		return
	}

	log.Printf("Successfully minted KYC NFT for wallet: %s", kycInfo.WalletAddress)

	// Return success response
	respondWithJSON(w, http.StatusOK, database.ApiResponse{
		Success: true,
		Message: "NFT minted successfully",
		Data: map[string]string{
			"wallet_address": kycInfo.WalletAddress,
			"citizen_id":     kycInfo.CitizenID,
		},
	})
}

// GetKYCByCitizenID handles retrieving KYC information by citizen ID
func (h *Handler) GetKYCByCitizenID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	citizenID := vars["citizen_id"]

	if citizenID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing citizen ID")
		return
	}

	kycInfo, err := h.repo.GetKYCByCitizenID(citizenID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "KYC information not found for this citizen ID")
		return
	}

	respondWithJSON(w, http.StatusOK, database.ApiResponse{
		Success: true,
		Message: "KYC information retrieved successfully",
		Data:    kycInfo,
	})
}

// GetKYCByWalletAddress handles retrieving KYC information by wallet address
func (h *Handler) GetKYCByWalletAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletAddress := vars["wallet_address"]

	if walletAddress == "" {
		respondWithError(w, http.StatusBadRequest, "Missing wallet address")
		return
	}

	kycInfo, err := h.repo.GetKYCByWalletAddress(walletAddress)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "KYC information not found for this wallet address")
		return
	}

	respondWithJSON(w, http.StatusOK, database.ApiResponse{
		Success: true,
		Message: "KYC information retrieved successfully",
		Data:    kycInfo,
	})
}

// respondWithError is a helper function to respond with an error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, database.ApiResponse{
		Success: false,
		Error:   message,
	})
}

// respondWithJSON is a helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
