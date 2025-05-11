package database

import (
	"time"
)

// Deposit định nghĩa cấu trúc cho bảng DEPOSITS
type Deposit struct {
	ID         int       `json:"id"`
	Commitment string    `json:"commitment"`
	Depositor  string    `json:"depositor"`
	LeafIndex  int       `json:"leaf_index"`
	Timestamp  time.Time `json:"timestamp"`
	TxHash     string    `json:"tx_hash"`
}

// Withdrawal định nghĩa cấu trúc cho bảng WITHDRAWALS
type Withdrawal struct {
	ID            int       `json:"id"`
	Recipient     string    `json:"recipient"`
	NullifierHash string    `json:"nullifier_hash"`
	Relayer       string    `json:"relayer"`
	Fee           float64   `json:"fee"`
	Timestamp     time.Time `json:"timestamp"`
	TxHash        string    `json:"tx_hash"`
}

// KYC định nghĩa cấu trúc cho bảng KYC
type KYC struct {
	CitizenID       string    `json:"citizen_id"`
	WalletAddress   string    `json:"wallet_address"`
	FullName        string    `json:"full_name"`
	PhoneNumber     string    `json:"phone_number"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Nationality     string    `json:"nationality"`
	KYCVerifiedAt   time.Time `json:"kyc_verified_at"`
	Verifier        string    `json:"verifier"`
	IsActive        bool      `json:"is_active"`
	WalletSignature string    `json:"wallet_signature"`
}
