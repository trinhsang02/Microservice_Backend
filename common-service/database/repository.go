package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"common-service/blockchain"
)

// NFTRecord represents an NFT minted for a KYC verification
type NFTRecord struct {
	ID            int        `json:"id"`
	WalletAddress string     `json:"wallet_address"`
	TokenID       string     `json:"token_id"`
	TokenURI      string     `json:"token_uri,omitempty"`
	CitizenID     string     `json:"citizen_id"`
	MintedAt      time.Time  `json:"minted_at"`
	TxHash        string     `json:"tx_hash,omitempty"`
	Status        string     `json:"status"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// Repository represents a data access layer for the database
type Repository struct {
	db interface{}
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		db: DB,
	}
}

// Submit KYC
func (r *Repository) SubmitKYC(kyc KYC) error {
	query := `
		INSERT INTO kyc (
			citizen_id, wallet_address, full_name, phone_number,
			date_of_birth, nationality, kyc_verified_at, verifier,
			is_active, wallet_signature
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) ON CONFLICT (citizen_id) 
		DO UPDATE SET 
			wallet_address = $2,
			full_name = $3,
			phone_number = $4,
			date_of_birth = $5,
			nationality = $6,
			kyc_verified_at = $7,
			verifier = $8,
			is_active = $9,
			wallet_signature = $10`

	_, err := DB.Exec(
		query,
		kyc.CitizenID, kyc.WalletAddress, kyc.FullName, kyc.PhoneNumber,
		kyc.DateOfBirth, kyc.Nationality, kyc.KYCVerifiedAt, kyc.Verifier,
		kyc.IsActive, kyc.WalletSignature,
	)

	if err != nil {
		log.Printf("Error submitting KYC: %v", err)
		return err
	}

	return nil
}

// KYCRecord is a struct used internally for scanning database results with NULL values
type KYCRecord struct {
	CitizenID       string
	WalletAddress   sql.NullString
	FullName        sql.NullString
	PhoneNumber     sql.NullString
	DateOfBirth     sql.NullTime
	Nationality     sql.NullString
	KYCVerifiedAt   sql.NullTime
	Verifier        sql.NullString
	IsActive        sql.NullBool
	WalletSignature sql.NullString
}

// ToKYC converts a KYCRecord to a KYC struct
func (r KYCRecord) ToKYC() KYC {
	kyc := KYC{
		CitizenID: r.CitizenID,
	}

	if r.WalletAddress.Valid {
		kyc.WalletAddress = r.WalletAddress.String
	}

	if r.FullName.Valid {
		kyc.FullName = r.FullName.String
	}

	if r.PhoneNumber.Valid {
		kyc.PhoneNumber = r.PhoneNumber.String
	}

	if r.DateOfBirth.Valid {
		kyc.DateOfBirth = r.DateOfBirth.Time
	}

	if r.Nationality.Valid {
		kyc.Nationality = r.Nationality.String
	}

	if r.KYCVerifiedAt.Valid {
		t := r.KYCVerifiedAt.Time
		kyc.KYCVerifiedAt = &t
	} else {
		// Default to current time if not valid
		now := time.Now()
		kyc.KYCVerifiedAt = &now
	}

	if r.Verifier.Valid {
		kyc.Verifier = r.Verifier.String
	} else {
		// Default to "system" if not valid
		kyc.Verifier = "system"
	}

	if r.IsActive.Valid {
		kyc.IsActive = r.IsActive.Bool
	}

	if r.WalletSignature.Valid {
		kyc.WalletSignature = r.WalletSignature.String
	}

	return kyc
}

// GetKYCByCitizenID retrieves KYC information by citizen ID
func (r *Repository) GetKYCByCitizenID(citizenID string) (*KYC, error) {
	log.Printf("Attempting to retrieve KYC by citizen ID: %s", citizenID)

	query := `
		SELECT 
			citizen_id, 
			wallet_address, 
			full_name, 
			phone_number,
			date_of_birth, 
			nationality, 
			kyc_verified_at, 
			verifier,
			is_active, 
			wallet_signature
		FROM kyc 
		WHERE citizen_id = $1`

	log.Printf("Executing query: %s with citizenID: %s", query, citizenID)

	var record KYCRecord
	err := DB.QueryRow(query, citizenID).Scan(
		&record.CitizenID,
		&record.WalletAddress,
		&record.FullName,
		&record.PhoneNumber,
		&record.DateOfBirth,
		&record.Nationality,
		&record.KYCVerifiedAt,
		&record.Verifier,
		&record.IsActive,
		&record.WalletSignature,
	)

	if err != nil {
		log.Printf("Error getting KYC by citizen ID: %v", err)
		return nil, err
	}

	kyc := record.ToKYC()
	log.Printf("Successfully retrieved KYC for citizen ID: %s", citizenID)
	return &kyc, nil
}

// GetKYCByWalletAddress retrieves KYC information by wallet address
func (r *Repository) GetKYCByWalletAddress(walletAddress string) (*KYC, error) {
	log.Printf("Attempting to retrieve KYC by wallet address: %s", walletAddress)

	query := `
		SELECT 
			citizen_id, 
			wallet_address, 
			full_name, 
			phone_number,
			date_of_birth, 
			nationality, 
			kyc_verified_at,
			verifier,
			kyc_verified_at, 
			verifier,
			is_active, 
			wallet_signature
		FROM kyc 
		WHERE wallet_address = $1`

	log.Printf("Executing query: %s with wallet address: %s", query, walletAddress)

	var record KYCRecord
	err := DB.QueryRow(query, walletAddress).Scan(
		&record.CitizenID,
		&record.WalletAddress,
		&record.FullName,
		&record.PhoneNumber,
		&record.DateOfBirth,
		&record.Nationality,
		&record.KYCVerifiedAt,
		&record.Verifier,
		&record.IsActive,
		&record.WalletSignature,
	)

	if err != nil {
		log.Printf("Error getting KYC by wallet address: %v", err)
		return nil, err
	}

	kyc := record.ToKYC()
	log.Printf("Successfully retrieved KYC for wallet address: %s", walletAddress)
	return &kyc, nil
}

// MintKYCNFT creates an NFT for a user who has completed KYC verification
func (r *Repository) MintKYCNFT(kyc KYC) error {
	log.Printf("Attempting to mint KYC NFT for wallet address: %s", kyc.WalletAddress)

	// Initialize blockchain service
	blockchainService, err := blockchain.NewNFTService()
	if err != nil {
		log.Printf("Error initializing blockchain service: %v", err)
		return fmt.Errorf("failed to initialize blockchain service: %v", err)
	}

	// Call blockchain service to mint the NFT
	mintResult, err := blockchainService.MintKYCNFT(kyc.WalletAddress, kyc.CitizenID)
	if err != nil {
		log.Printf("Error minting NFT on blockchain: %v", err)
		return fmt.Errorf("failed to mint NFT on blockchain: %v", err)
	}

	log.Printf("Successfully minted KYC NFT with token ID: %s for wallet: %s, TxHash: %s",
		mintResult.TokenID, kyc.WalletAddress, mintResult.TxHash)

	return nil
}
