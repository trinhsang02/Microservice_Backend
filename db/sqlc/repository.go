package sqlc

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// Repository holds the database queries.
type Repository struct {
	queries *Queries
}

// NewRepository creates a new Repository instance.
func NewRepository(queries *Queries) *Repository {
	return &Repository{queries: queries}
}

// SubmitKYC inserts a new KYC record and associated wallet info.
func (r *Repository) SubmitKYC(ctx context.Context, kyc KycInfo, walletAddress string, walletSignature string) error {
	// Create KYC record
	_, err := r.queries.CreateKycInfo(ctx, CreateKycInfoParams{
		CitizenID:     kyc.CitizenID,
		FullName:      kyc.FullName,
		PhoneNumber:   kyc.PhoneNumber,
		DateOfBirth:   kyc.DateOfBirth,
		Nationality:   kyc.Nationality,
		Verifier:      kyc.Verifier,
		IsActive:      kyc.IsActive,
		KycVerifiedAt: kyc.KycVerifiedAt,
	})
	if err != nil {
		return err
	}

	// Create or update wallet info
	err = r.queries.CreateOrUpdateWalletInfo(ctx, CreateOrUpdateWalletInfoParams{
		WalletAddress:   walletAddress,
		CitizenID:       pgtype.Text{String: kyc.CitizenID, Valid: true},
		WalletSignature: pgtype.Text{String: walletSignature, Valid: true},
	})
	return err
}

// GetKYCByCitizenID retrieves KYC info by citizen ID.
func (r *Repository) GetKYCByCitizenID(ctx context.Context, citizenID string) (*KycInfo, error) {
	kyc, err := r.queries.GetKycInfoByCitizenID(ctx, citizenID)
	if err != nil {
		return nil, err
	}
	return &KycInfo{
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
func (r *Repository) GetKYCByWalletAddress(ctx context.Context, walletAddress string) (*KycInfo, error) {
	kyc, err := r.queries.GetKycInfoByWalletAddress(ctx, walletAddress)
	if err != nil {
		return nil, err
	}
	return &KycInfo{
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
func (r *Repository) UpdateKYC(ctx context.Context, kyc KycInfo) error {
	_, err := r.queries.UpdateKycInfo(ctx, UpdateKycInfoParams{
		CitizenID:     kyc.CitizenID,
		FullName:      kyc.FullName,
		PhoneNumber:   kyc.PhoneNumber,
		DateOfBirth:   kyc.DateOfBirth,
		Nationality:   kyc.Nationality,
		Verifier:      kyc.Verifier,
		IsActive:      kyc.IsActive,
		KycVerifiedAt: kyc.KycVerifiedAt,
	})
	return err
}

func (r *Repository) GetDepositEventsFromBlockToBlock(ctx context.Context, netId string, contractAddress string, fromBlock string, toBlock string) ([]Deposit, error) {
	fromBlockInt, _ := strconv.ParseInt(fromBlock, 10, 32)
	toBlockInt, _ := strconv.ParseInt(toBlock, 10, 32)
	netIdInt, _ := strconv.ParseInt(netId, 10, 32)

	if fromBlockInt == 0 {
		res, err := r.queries.GetEarliestDepositSyncedBlock(ctx, GetEarliestDepositSyncedBlockParams{
			ContractAddress: pgtype.Text{String: contractAddress, Valid: true},
			ChainID:         pgtype.Int4{Int32: int32(netIdInt), Valid: true},
		})
		if err == nil {
			fromBlockInt = int64(res.(int32))
		}
	}

	if toBlockInt == 0 {
		res, err := r.queries.GetLatestDepositSyncedBlock(ctx, GetLatestDepositSyncedBlockParams{
			ContractAddress: pgtype.Text{String: contractAddress, Valid: true},
			ChainID:         pgtype.Int4{Int32: int32(netIdInt), Valid: true},
		})
		if err == nil {
			toBlockInt = int64(res.(int32))
		}
	}

	events, err := r.queries.GetDepositsFromBlockToBlock(ctx, GetDepositsFromBlockToBlockParams{
		BlockNumber:     pgtype.Int4{Int32: int32(fromBlockInt), Valid: true},
		BlockNumber_2:   pgtype.Int4{Int32: int32(toBlockInt), Valid: true},
		ChainID:         pgtype.Int4{Int32: int32(netIdInt), Valid: true},
		ContractAddress: pgtype.Text{String: contractAddress, Valid: true},
	})
	return events, err
}

func (r *Repository) GetWithdrawalByNullifierHash(ctx context.Context, nullifierHash string) (*Withdrawal, error) {

	withdrawal, err := r.queries.GetWithdrawalByNullifierHash(ctx, pgtype.Text{String: nullifierHash, Valid: true})

	return &Withdrawal{
		ID:              withdrawal.ID,
		NullifierHash:   withdrawal.NullifierHash,
		ContractAddress: withdrawal.ContractAddress,
		ChainID:         withdrawal.ChainID,
		TxHash:          withdrawal.TxHash,
		Timestamp:       withdrawal.Timestamp,
		BlockNumber:     withdrawal.BlockNumber,
	}, err
}

func (r *Repository) GetDepositByCommitment(ctx context.Context, commitment string) (*Deposit, error) {

	deposit, err := r.queries.GetDepositByCommitment(ctx, pgtype.Text{String: commitment, Valid: true})

	return &Deposit{
		ID:              deposit.ID,
		ContractAddress: deposit.ContractAddress,
		Commitment:      deposit.Commitment,
		Depositor:       deposit.Depositor,
		LeafIndex:       deposit.LeafIndex,
		TxHash:          deposit.TxHash,
		Timestamp:       deposit.Timestamp,
		BlockNumber:     deposit.BlockNumber,
		ChainID:         deposit.ChainID,
	}, err
}

func (r *Repository) GetLeaves(ctx context.Context, netId string, contractAddress string) ([]string, error) {
	netIdInt, _ := strconv.ParseInt(netId, 10, 32)

	leaves, err := r.queries.GetLeaves(ctx, GetLeavesParams{
		ChainID:         pgtype.Int4{Int32: int32(netIdInt), Valid: true},
		ContractAddress: pgtype.Text{String: contractAddress, Valid: true},
	})

	result := make([]string, len(leaves))
	for i, leaf := range leaves {
		result[i] = leaf.String
	}
	return result, err
}

// GetKYCStatusByWalletAddress returns only the is_active status for a wallet address.
func (r *Repository) GetKYCStatusByWalletAddress(ctx context.Context, walletAddress string) (pgtype.Bool, error) {
	return r.queries.GetKycStatusByWalletAddress(ctx, walletAddress)
}
