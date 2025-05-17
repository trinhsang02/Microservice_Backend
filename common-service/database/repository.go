package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourusername/yourrepo/db/sqlc"
)

// Repository holds the database queries.
type Repository struct {
	queries *sqlc.Queries
}

// NewRepository creates a new Repository instance.
func NewRepository(queries *sqlc.Queries) *Repository {
	return &Repository{queries: queries}
}

// SubmitKYC inserts a new KYC record.
func (r *Repository) SubmitKYC(ctx context.Context, kyc sqlc.KycInfo) error {
	_, err := r.queries.CreateKycInfo(ctx, sqlc.CreateKycInfoParams{
		CitizenID:     kyc.CitizenID,
		FullName:      pgtype.Text{String: kyc.FullName.String, Valid: true},
		PhoneNumber:   pgtype.Text{String: kyc.PhoneNumber.String, Valid: true},
		DateOfBirth:   pgtype.Date{Time: kyc.DateOfBirth.Time, Valid: true},
		Nationality:   pgtype.Text{String: kyc.Nationality.String, Valid: true},
		Verifier:      pgtype.Text{String: kyc.Verifier.String, Valid: true},
		IsActive:      pgtype.Bool{Bool: kyc.IsActive.Bool, Valid: true},
		KycVerifiedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	return err
}

// GetKYCByCitizenID retrieves KYC info by citizen ID.
func (r *Repository) GetKYCByCitizenID(ctx context.Context, citizenID string) (*sqlc.KycInfo, error) {
	kyc, err := r.queries.GetKycInfoByCitizenID(ctx, citizenID)
	if err != nil {
		return nil, err
	}
	return &sqlc.KycInfo{
		CitizenID:     kyc.CitizenID,
		FullName:      pgtype.Text{String: kyc.FullName.String, Valid: true},
		PhoneNumber:   pgtype.Text{String: kyc.PhoneNumber.String, Valid: true},
		DateOfBirth:   pgtype.Date{Time: kyc.DateOfBirth.Time, Valid: true},
		Nationality:   pgtype.Text{String: kyc.Nationality.String, Valid: true},
		Verifier:      pgtype.Text{String: kyc.Verifier.String, Valid: true},
		IsActive:      pgtype.Bool{Bool: kyc.IsActive.Bool, Valid: true},
		KycVerifiedAt: pgtype.Timestamp{Time: kyc.KycVerifiedAt.Time, Valid: true},
	}, nil
}

// GetKYCByWalletAddress retrieves KYC info by wallet address.
func (r *Repository) GetKYCByWalletAddress(ctx context.Context, walletAddress string) (*sqlc.KycInfo, error) {
	kyc, err := r.queries.GetKycInfoByWalletAddress(ctx, walletAddress)
	if err != nil {
		return nil, err
	}
	return &sqlc.KycInfo{
		CitizenID:     kyc.CitizenID,
		FullName:      pgtype.Text{String: kyc.FullName.String, Valid: true},
		PhoneNumber:   pgtype.Text{String: kyc.PhoneNumber.String, Valid: true},
		DateOfBirth:   pgtype.Date{Time: kyc.DateOfBirth.Time, Valid: true},
		Nationality:   pgtype.Text{String: kyc.Nationality.String, Valid: true},
		Verifier:      pgtype.Text{String: kyc.Verifier.String, Valid: true},
		IsActive:      pgtype.Bool{Bool: kyc.IsActive.Bool, Valid: true},
		KycVerifiedAt: pgtype.Timestamp{Time: kyc.KycVerifiedAt.Time, Valid: true},
	}, nil
}

// UpdateKYC updates an existing KYC record.
func (r *Repository) UpdateKYC(ctx context.Context, kyc sqlc.KycInfo) error {
	_, err := r.queries.UpdateKycInfo(ctx, sqlc.UpdateKycInfoParams{
		CitizenID:     kyc.CitizenID,
		FullName:      pgtype.Text{String: kyc.FullName.String, Valid: true},
		PhoneNumber:   pgtype.Text{String: kyc.PhoneNumber.String, Valid: true},
		DateOfBirth:   pgtype.Date{Time: kyc.DateOfBirth.Time, Valid: true},
		Nationality:   pgtype.Text{String: kyc.Nationality.String, Valid: true},
		Verifier:      pgtype.Text{String: kyc.Verifier.String, Valid: true},
		IsActive:      pgtype.Bool{Bool: kyc.IsActive.Bool, Valid: true},
		KycVerifiedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	return err
}
