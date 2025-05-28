package api

import (
	"log"
	"net/http"
	// "os"
	// "path/filepath"
	// "github.com/xuri/excelize/v2"
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

// logKYCToExcel logs KYC information to an Excel file
// func logKYCToExcel(kyc sqlc.KycInfo, walletAddress string, status string) error {
// 	logsDir := "logs"
// 	log.Printf("Creating logs directory at: %s", logsDir)
// 	if err := os.MkdirAll(logsDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create logs directory: %v", err)
// 	}

// 	currentDate := time.Now().Format("2006-01-02")
// 	filename := filepath.Join(logsDir, fmt.Sprintf("kyc_logs_%s.xlsx", currentDate))
// 	log.Printf("Excel file path: %s", filename)

// 	var f *excelize.File
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		log.Printf("Creating new Excel file: %s", filename)
// 		f = excelize.NewFile()
// 		headers := []string{"Timestamp", "Citizen ID", "Full Name", "Phone Number", "Date of Birth",
// 			"Nationality", "Verifier", "Wallet Address", "Status"}
// 		for i, header := range headers {
// 			cell := fmt.Sprintf("%c1", 'A'+i)
// 			f.SetCellValue("Sheet1", cell, header)
// 		}
// 	} else {
// 		log.Printf("Opening existing Excel file: %s", filename)
// 		var err error
// 		f, err = excelize.OpenFile(filename)
// 		if err != nil {
// 			return fmt.Errorf("failed to open existing file: %v", err)
// 		}
// 	}

// 	rows, err := f.GetRows("Sheet1")
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows: %v", err)
// 	}
// 	lastRow := len(rows) + 1
// 	log.Printf("Adding new row at position: %d", lastRow)

// 	row := []interface{}{
// 		time.Now().Format("2006-01-02 15:04:05"),
// 		kyc.CitizenID,
// 		kyc.FullName.String,
// 		kyc.PhoneNumber.String,
// 		kyc.DateOfBirth.Time.Format("2006-01-02"),
// 		kyc.Nationality.String,
// 		kyc.Verifier.String,
// 		walletAddress,
// 		status,
// 	}

// 	for i, value := range row {
// 		cell := fmt.Sprintf("%c%d", 'A'+i, lastRow)
// 		f.SetCellValue("Sheet1", cell, value)
// 	}

// 	log.Printf("Saving Excel file: %s", filename)
// 	if err := f.SaveAs(filename); err != nil {
// 		return fmt.Errorf("failed to save file: %v", err)
// 	}
// 	log.Printf("Successfully saved Excel file")
// 	return nil
// }

// SubmitKYC handles the submission of KYC information.
func (h *Handler) SubmitKYC(c *gin.Context) {
	var req KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received KYC submission request for wallet: %s", req.WalletAddress)

	// Check if KYC already exists for this wallet
	existingKYC, err := h.repo.GetKYCByWalletAddress(c.Request.Context(), req.WalletAddress)
	if err == nil && existingKYC != nil {
		// Check if KYC is already active
		if existingKYC.IsActive.Bool {
			log.Printf("KYC is already active for wallet %s", req.WalletAddress)
			c.JSON(http.StatusBadRequest, gin.H{"error": "KYC is already active for this wallet"})
			return
		}

		// KYC exists but not active, proceed with mint
		log.Printf("KYC exists but not active for wallet %s, proceeding with mint", req.WalletAddress)

		// Log to Excel
		// if err := logKYCToExcel(*existingKYC, req.WalletAddress, "MINT_REQUESTED"); err != nil {
		// 	log.Printf("Failed to log KYC to Excel: %v", err)
		// }

		mintMsg := MintMessage{
			KycInfo:       *existingKYC,
			WalletAddress: req.WalletAddress,
		}

		log.Printf("Publishing mint message for wallet: %s", req.WalletAddress)
		if err := h.producer.PublishStruct("kyc.mint", mintMsg); err != nil {
			log.Printf("Failed to publish message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process KYC submission"})
			return
		}

		log.Println("Mint message published successfully")
		c.JSON(http.StatusOK, gin.H{"message": "KYC exists but not active, proceeding with mint"})
		return
	}

	// If KYC doesn't exist, create new KYC
	kyc := sqlc.KycInfo{
		CitizenID:   req.CitizenID,
		FullName:    pgtype.Text{String: req.FullName, Valid: true},
		PhoneNumber: pgtype.Text{String: req.PhoneNumber, Valid: true},
		DateOfBirth: pgtype.Date{Time: time.Now(), Valid: true}, // Parse req.DateOfBirth string to time.Time
		Nationality: pgtype.Text{String: req.Nationality, Valid: true},
		Verifier:    pgtype.Text{String: req.Verifier, Valid: true},
		IsActive:    pgtype.Bool{Bool: false, Valid: true},
	}

	if err := h.repo.SubmitKYC(c.Request.Context(), kyc, req.WalletAddress, req.WalletSignature); err != nil {
		log.Printf("Failed to submit KYC to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log to Excel
	// if err := logKYCToExcel(kyc, req.WalletAddress, "SUBMITTED"); err != nil {
	// 	log.Printf("Failed to log KYC to Excel: %v", err)
	// }

	log.Println("KYC submitted to database successfully")

	// Prepare message for RabbitMQ
	mintMsg := MintMessage{
		KycInfo:       kyc,
		WalletAddress: req.WalletAddress,
	}

	log.Printf("Publishing mint message for wallet: %s", req.WalletAddress)
	if err := h.producer.PublishStruct("kyc.mint", mintMsg); err != nil {
		log.Printf("Failed to publish message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process KYC submission"})
		return
	}

	log.Println("Mint message published successfully")
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