package database

import (
	"time"
)

// ApiResponse is a standard response format for API calls
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// KYC represents user KYC information according to the new schema
type KYC struct {
	CitizenID       string     `json:"citizen_id"`
	WalletAddress   string     `json:"wallet_address,omitempty"`
	FullName        string     `json:"full_name,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty"`
	DateOfBirth     time.Time  `json:"date_of_birth,omitempty"`
	Nationality     string     `json:"nationality,omitempty"`
	KYCVerifiedAt   *time.Time `json:"kyc_verified_at"`
	Verifier        string     `json:"verifier"`
	IsActive        bool       `json:"is_active,omitempty"`
	WalletSignature string     `json:"wallet_signature,omitempty"`
}

// NFTMintEvent represents the event for minting an NFT after KYC
type NFTMintEvent struct {
	WalletAddress string    `json:"wallet_address"`
	CitizenID     string    `json:"citizen_id"`
	FullName      string    `json:"full_name"`
	Timestamp     time.Time `json:"timestamp"`
	EventType     string    `json:"event_type"`
	Status        string    `json:"status,omitempty"`
	TxHash        string    `json:"tx_hash,omitempty"`
	TokenID       string    `json:"token_id,omitempty"`
}
