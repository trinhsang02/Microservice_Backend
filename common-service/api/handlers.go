package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourusername/yourrepo/db/sqlc"
)

// Handler struct holds dependencies for API handlers
type Handler struct {
	repo *sqlc.Repository
}

// NewHandler creates a new Handler instance
func NewHandler(repo *sqlc.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// KYCRequest represents the request payload for KYC submission
type KYCRequest struct {
	CitizenID     string `json:"citizen_id"`
	FullName      string `json:"full_name"`
	PhoneNumber   string `json:"phone_number"`
	DateOfBirth   string `json:"date_of_birth"` // Format: YYYY-MM-DD
	Nationality   string `json:"nationality"`
	Verifier      string `json:"verifier,omitempty"`
	IsActive      bool   `json:"is_active"`
	KYCVerifiedAt string `json:"kyc_verified_at,omitempty"` // Format: YYYY-MM-DD HH:MM:SS or empty
}

// SubmitKYC handles the submission of KYC information.
func (h *Handler) SubmitKYC(c *gin.Context) {
	var req KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kyc := sqlc.KycInfo{
		CitizenID:   req.CitizenID,
		FullName:    pgtype.Text{String: req.FullName, Valid: true},
		PhoneNumber: pgtype.Text{String: req.PhoneNumber, Valid: true},
		DateOfBirth: pgtype.Date{Time: time.Now(), Valid: true}, // Parse req.DateOfBirth string to time.Time
		Nationality: pgtype.Text{String: req.Nationality, Valid: true},
		Verifier:    pgtype.Text{String: req.Verifier, Valid: true},
		IsActive:    pgtype.Bool{Bool: true, Valid: true},
	}

	if err := h.repo.SubmitKYC(c.Request.Context(), kyc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC submitted successfully"})
}

// GetKYCByCitizenID retrieves KYC information by citizen ID.
func (h *Handler) GetKYCByCitizenID(c *gin.Context) {
	citizenID := c.Param("citizenID")
	kyc, err := h.repo.GetKYCByCitizenID(c.Request.Context(), citizenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kyc)
}

// GetKYCByWalletAddress retrieves KYC information by wallet address.
func (h *Handler) GetKYCByWalletAddress(c *gin.Context) {
	walletAddress := c.Param("walletAddress")
	kyc, err := h.repo.GetKYCByWalletAddress(c.Request.Context(), walletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kyc)
}

// UpdateKYC updates KYC information.
func (h *Handler) UpdateKYC(c *gin.Context) {
	var req KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kyc := sqlc.KycInfo{
		CitizenID:   req.CitizenID,
		FullName:    pgtype.Text{String: req.FullName, Valid: true},
		PhoneNumber: pgtype.Text{String: req.PhoneNumber, Valid: true},
		DateOfBirth: pgtype.Date{Time: time.Now(), Valid: true}, // Parse req.DateOfBirth string to time.Time
		Nationality: pgtype.Text{String: req.Nationality, Valid: true},
		Verifier:    pgtype.Text{String: req.Verifier, Valid: true},
		IsActive:    pgtype.Bool{Bool: true, Valid: true},
	}

	if err := h.repo.UpdateKYC(c.Request.Context(), kyc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC updated successfully"})
}

// GetEvents handles retrieving events from a specific block
func (h *Handler) GetEvents(c *gin.Context) {
	netId := c.Param("netId")
	contractAddress := c.Param("contractAddress")
	eventType := c.Param("eventType")
	fromBlock := c.DefaultQuery("fromBlock", "0")
	toBlock := c.DefaultQuery("toBlock", "0")
	limit := c.DefaultQuery("limit", "0")

	// Validate parameters
	if netId == "" || contractAddress == "" || eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters",
		})
		return
	}

	var events []sqlc.Deposit
	var err error

	if strings.ToLower(eventType) == "withdrawal" {
		// TODO: Implement withdrawal event retrieval
	} else if strings.ToLower(eventType) == "deposit" {
		events, err = h.repo.GetDepositEventsFromBlockToBlock(c.Request.Context(), netId, contractAddress, fromBlock, toBlock)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch events",
		})
		return
	}
	// Check if we need to limit the number of events returned
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0 // Default to 0 (no limit) if conversion fails
	}

	var responseEvents []sqlc.Deposit
	if limitInt > 0 && limitInt < len(events) {
		responseEvents = events[:limitInt]
	} else {
		responseEvents = events
	}

	c.JSON(http.StatusOK, gin.H{
		"events": responseEvents,
		"count":  len(responseEvents),
	})

}

func (h *Handler) GetEventByInfo(c *gin.Context) {
	eventType := c.Param("eventType")
	hex := c.Param("hex")

	// Validate parameters
	if eventType == "" || hex == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters",
		})
		return
	}

	if strings.ToLower(eventType) == "withdrawal" {
		// TODO: Implement withdrawal event retrieval
		event, err := h.repo.GetWithdrawalByNullifierHash(c.Request.Context(), hex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch event",
			})
			log.Println("error", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"event": event,
		})
	} else if strings.ToLower(eventType) == "deposit" {
		// TODO: Implement deposit event retrieval
		event, err := h.repo.GetDepositByCommitment(c.Request.Context(), hex)
		log.Println("Found deposit event", event)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch event",
			})
			log.Println("error", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"event": event,
		})
	}
}

func (h *Handler) GetLeaves(c *gin.Context) {
	netId := c.Param("netId")
	contractAddress := c.Param("contractAddress")

	leaves, err := h.repo.GetLeaves(c.Request.Context(), netId, contractAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch leaves",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

// MintKYCNFT handles minting an NFT for a user with verified KYC
// func (h *Handler) MintKYCNFT(w http.ResponseWriter, r *http.Request) {
// 	var request struct {
// 		WalletAddress string `json:"wallet_address"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&request); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	defer r.Body.Close()

// 	if request.WalletAddress == "" {
// 		respondWithError(w, http.StatusBadRequest, "Wallet address is required")
// 		return
// 	}

// 	// Get KYC information for the wallet address
// 	kycInfo, err := h.repo.GetKYCByWalletAddress(request.WalletAddress)
// 	if err != nil {
// 		respondWithError(w, http.StatusNotFound, "KYC information not found for this wallet address")
// 		return
// 	}

// 	// Verify that KYC is active
// 	if !kycInfo.IsActive {
// 		respondWithError(w, http.StatusBadRequest, "KYC is not active for this wallet address")
// 		return
// 	}

// 	// Mint NFT
// 	err = h.repo.MintKYCNFT(*kycInfo)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Failed to mint NFT: "+err.Error())
// 		return
// 	}

// 	log.Printf("Successfully minted KYC NFT for wallet: %s", kycInfo.WalletAddress)

// 	// Return success response
// 	respondWithJSON(w, http.StatusOK, database.ApiResponse{
// 		Success: true,
// 		Message: "NFT minted successfully",
// 		Data: map[string]string{
// 			"wallet_address": kycInfo.WalletAddress,
// 			"citizen_id":     kycInfo.CitizenID,
// 		},
// 	})
// }

// // respondWithError is a helper function to respond with an error
// func respondWithError(w http.ResponseWriter, code int, message string) {
// 	respondWithJSON(w, code, database.ApiResponse{
// 		Success: false,
// 		Error:   message,
// 	})
// }

// // respondWithJSON is a helper function to respond with JSON
// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }
