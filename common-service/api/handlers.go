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
	"github.com/yourusername/yourrepo/mq/rabbitmq"
)

// Handler struct holds dependencies for API handlers
type Handler struct {
	repo     *sqlc.Repository
	producer *rabbitmq.Producer
}

// NewHandler creates a new Handler instance
func NewHandler(repo *sqlc.Repository, producer *rabbitmq.Producer) *Handler {
	return &Handler{
		repo:     repo,
		producer: producer,
	}
}

// KYCRequest represents the request payload for KYC submission
type KYCRequest struct {
	CitizenID       string `json:"citizen_id"`
	FullName        string `json:"full_name"`
	PhoneNumber     string `json:"phone_number"`
	DateOfBirth     string `json:"date_of_birth"` // Format: YYYY-MM-DD
	Nationality     string `json:"nationality"`
	Verifier        string `json:"verifier,omitempty"`
	IsActive        bool   `json:"is_active"`
	KYCVerifiedAt   string `json:"kyc_verified_at,omitempty"` // Format: YYYY-MM-DD HH:MM:SS or empty
	WalletAddress   string `json:"wallet_address"`
	WalletSignature string `json:"wallet_signature"`
}

// MintMessage represents the message structure for minting NFT
type MintMessage struct {
	sqlc.KycInfo  `json:"kyc"`
	WalletAddress string `json:"wallet_address"`
}

// SubmitKYC handles the submission of KYC information.
func (h *Handler) SubmitKYC(c *gin.Context) {
	var req KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.CitizenID == "" || req.FullName == "" || req.WalletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Debug log: Print request data
	log.Printf("KYC Request Data: CitizenID=%s, FullName=%s, PhoneNumber=%s, DateOfBirth=%s, Nationality=%s, Verifier=%s, WalletAddress=%s",
		req.CitizenID, req.FullName, req.PhoneNumber, req.DateOfBirth, req.Nationality, req.Verifier, req.WalletAddress)

	// Check if KYC already exists for this wallet
	existingKYC, err := h.repo.GetKYCByWalletAddress(c.Request.Context(), req.WalletAddress)
	if err == nil && existingKYC != nil {
		if existingKYC.IsActive.Bool {
			c.JSON(http.StatusConflict, gin.H{"error": "KYC is already active for this wallet"})
			return
		}
		mintMsg := MintMessage{
			KycInfo:       *existingKYC,
			WalletAddress: req.WalletAddress,
		}
		if err := h.producer.PublishStruct("kyc.mint", mintMsg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process KYC submission"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "KYC exists but not active, proceeding with mint"})
		return
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
		return
	}

	// Parse KYC verified at if provided
	var kycVerifiedAt pgtype.Timestamp
	if req.KYCVerifiedAt != "" {
		verifiedAt, err := time.Parse("2006-01-02 15:04:05", req.KYCVerifiedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kyc_verified_at format. Use YYYY-MM-DD HH:MM:SS"})
			return
		}
		kycVerifiedAt = pgtype.Timestamp{Time: verifiedAt, Valid: true}
	} else {
		kycVerifiedAt = pgtype.Timestamp{Time: time.Now(), Valid: true}
	}

	kyc := sqlc.KycInfo{
		CitizenID:     req.CitizenID,
		FullName:      pgtype.Text{String: req.FullName, Valid: true},
		PhoneNumber:   pgtype.Text{String: req.PhoneNumber, Valid: true},
		DateOfBirth:   pgtype.Date{Time: dob, Valid: true},
		Nationality:   pgtype.Text{String: req.Nationality, Valid: true},
		Verifier:      pgtype.Text{String: req.Verifier, Valid: req.Verifier != ""},
		IsActive:      pgtype.Bool{Bool: false, Valid: true}, // Khi tạo mới luôn là false
		KycVerifiedAt: kycVerifiedAt,
	}

	// Debug log: Print created KYC struct
	log.Printf("Created KYC Struct: CitizenID=%s, FullName={String:%s, Valid:%t}, PhoneNumber={String:%s, Valid:%t}, Nationality={String:%s, Valid:%t}",
		kyc.CitizenID, kyc.FullName.String, kyc.FullName.Valid, kyc.PhoneNumber.String, kyc.PhoneNumber.Valid, kyc.Nationality.String, kyc.Nationality.Valid)

	if err := h.repo.SubmitKYC(c.Request.Context(), kyc, req.WalletAddress, req.WalletSignature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save KYC information"})
		return
	}

	mintMsg := MintMessage{
		KycInfo:       kyc,
		WalletAddress: req.WalletAddress,
	}
	if err := h.producer.PublishStruct("kyc.mint", mintMsg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process KYC submission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "KYC submitted successfully",
		"data": gin.H{
			"citizen_id":     kyc.CitizenID,
			"wallet_address": req.WalletAddress,
			"is_active":      kyc.IsActive.Bool,
		},
	})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
		return
	}

	var kycVerifiedAt pgtype.Timestamp
	if req.KYCVerifiedAt != "" {
		verifiedAt, err := time.Parse("2006-01-02 15:04:05", req.KYCVerifiedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid kyc_verified_at format. Use YYYY-MM-DD HH:MM:SS"})
			return
		}
		kycVerifiedAt = pgtype.Timestamp{Time: verifiedAt, Valid: true}
	} else {
		kycVerifiedAt = pgtype.Timestamp{Time: time.Now(), Valid: true}
	}

	kyc := sqlc.KycInfo{
		CitizenID:     req.CitizenID,
		FullName:      pgtype.Text{String: req.FullName, Valid: true},
		PhoneNumber:   pgtype.Text{String: req.PhoneNumber, Valid: true},
		DateOfBirth:   pgtype.Date{Time: dob, Valid: true},
		Nationality:   pgtype.Text{String: req.Nationality, Valid: true},
		Verifier:      pgtype.Text{String: req.Verifier, Valid: req.Verifier != ""},
		IsActive:      pgtype.Bool{Bool: req.IsActive, Valid: true},
		KycVerifiedAt: kycVerifiedAt,
	}

	if err := h.repo.UpdateKYC(c.Request.Context(), kyc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KYC information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC updated successfully"})
}

// CheckKYCStatusByWalletAddress checks if KYC is active by wallet address.
func (h *Handler) CheckKYCStatusByWalletAddress(c *gin.Context) {
	walletAddress := c.Param("walletAddress")
	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing walletAddress parameter"})
		return
	}

	isActive, err := h.repo.GetKYCStatusByWalletAddress(c.Request.Context(), walletAddress)
	if err != nil || !isActive.Valid {
		c.JSON(http.StatusNotFound, gin.H{
			"error":          "KYC not found",
			"wallet_address": walletAddress,
			"is_active":      false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet_address": walletAddress,
		"is_active":      isActive.Bool,
	})
}

type EventUriParams struct {
	NetId           string `uri:"netId" binding:"required"`
	ContractAddress string `uri:"contractAddress" binding:"required"`
	EventType       string `uri:"eventType" binding:"required,oneof=withdrawal deposit"`
}

type EventQueryParams struct {
	FromBlock string `form:"fromBlock" binding:"required"`
	ToBlock   string `form:"toBlock" binding:"required"`
	Limit     string `form:"limit" binding:"required"`
}

// GetEvents handles retrieving events from a specific block
func (h *Handler) GetEvents(c *gin.Context) {
	var uriParams EventUriParams
	var queryParams EventQueryParams

	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var events []sqlc.Deposit
	var err error

	if strings.ToLower(uriParams.EventType) == "withdrawal" {
		// TODO: Implement withdrawal event retrieval
	} else if strings.ToLower(uriParams.EventType) == "deposit" {
		events, err = h.repo.GetDepositEventsFromBlockToBlock(c.Request.Context(), uriParams.NetId, uriParams.ContractAddress, queryParams.FromBlock, queryParams.ToBlock)
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
	limitInt, err := strconv.Atoi(queryParams.Limit)
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

type GetEventByInfoUriParams struct {
	EventType string `uri:"eventType" binding:"required,oneof=withdrawal deposit"`
	Hex       string `uri:"hex" binding:"required"`
}

func (h *Handler) GetEventByInfo(c *gin.Context) {
	var uriParams GetEventByInfoUriParams
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(uriParams.EventType) == "withdrawal" {
		event, err := h.repo.GetWithdrawalByNullifierHash(c.Request.Context(), uriParams.Hex)
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
	} else if strings.ToLower(uriParams.EventType) == "deposit" {
		event, err := h.repo.GetDepositByCommitment(c.Request.Context(), uriParams.Hex)
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

type GetLeavesUriParams struct {
	NetId           string `uri:"netId" binding:"required"`
	ContractAddress string `uri:"contractAddress" binding:"required"`
}

func (h *Handler) GetLeaves(c *gin.Context) {
	var uriParams GetLeavesUriParams
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leaves, err := h.repo.GetLeaves(c.Request.Context(), uriParams.NetId, uriParams.ContractAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch leaves",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}
