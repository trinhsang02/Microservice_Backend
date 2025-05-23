// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: withdrawals.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWithdrawal = `-- name: CreateWithdrawal :one
INSERT INTO withdrawals (contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (nullifier_hash) DO NOTHING
RETURNING id, contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id
`

type CreateWithdrawalParams struct {
	ContractAddress pgtype.Text
	Recipient       pgtype.Text
	NullifierHash   pgtype.Text
	Relayer         pgtype.Text
	Fee             pgtype.Numeric
	TxHash          pgtype.Text
	Timestamp       pgtype.Numeric
	BlockNumber     pgtype.Int4
	ChainID         pgtype.Int4
}

func (q *Queries) CreateWithdrawal(ctx context.Context, arg CreateWithdrawalParams) (Withdrawal, error) {
	row := q.db.QueryRow(ctx, createWithdrawal,
		arg.ContractAddress,
		arg.Recipient,
		arg.NullifierHash,
		arg.Relayer,
		arg.Fee,
		arg.TxHash,
		arg.Timestamp,
		arg.BlockNumber,
		arg.ChainID,
	)
	var i Withdrawal
	err := row.Scan(
		&i.ID,
		&i.ContractAddress,
		&i.Recipient,
		&i.NullifierHash,
		&i.Relayer,
		&i.Fee,
		&i.TxHash,
		&i.Timestamp,
		&i.BlockNumber,
		&i.ChainID,
	)
	return i, err
}

const getAllWithdrawalsOfContract = `-- name: GetAllWithdrawalsOfContract :many
SELECT id, contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id FROM withdrawals WHERE contract_address = $1
ORDER BY timestamp ASC
`

func (q *Queries) GetAllWithdrawalsOfContract(ctx context.Context, contractAddress pgtype.Text) ([]Withdrawal, error) {
	rows, err := q.db.Query(ctx, getAllWithdrawalsOfContract, contractAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Withdrawal
	for rows.Next() {
		var i Withdrawal
		if err := rows.Scan(
			&i.ID,
			&i.ContractAddress,
			&i.Recipient,
			&i.NullifierHash,
			&i.Relayer,
			&i.Fee,
			&i.TxHash,
			&i.Timestamp,
			&i.BlockNumber,
			&i.ChainID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllWithdrawalsOfRecipient = `-- name: GetAllWithdrawalsOfRecipient :many
SELECT id, contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id FROM withdrawals WHERE recipient = $1
ORDER BY timestamp ASC
`

func (q *Queries) GetAllWithdrawalsOfRecipient(ctx context.Context, recipient pgtype.Text) ([]Withdrawal, error) {
	rows, err := q.db.Query(ctx, getAllWithdrawalsOfRecipient, recipient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Withdrawal
	for rows.Next() {
		var i Withdrawal
		if err := rows.Scan(
			&i.ID,
			&i.ContractAddress,
			&i.Recipient,
			&i.NullifierHash,
			&i.Relayer,
			&i.Fee,
			&i.TxHash,
			&i.Timestamp,
			&i.BlockNumber,
			&i.ChainID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLatestWithdrawalSyncedBlockOfContractOnChain = `-- name: GetLatestWithdrawalSyncedBlockOfContractOnChain :one
SELECT MAX(block_number) FROM withdrawals WHERE contract_address = $1 AND chain_id = $2
`

type GetLatestWithdrawalSyncedBlockOfContractOnChainParams struct {
	ContractAddress pgtype.Text
	ChainID         pgtype.Int4
}

func (q *Queries) GetLatestWithdrawalSyncedBlockOfContractOnChain(ctx context.Context, arg GetLatestWithdrawalSyncedBlockOfContractOnChainParams) (interface{}, error) {
	row := q.db.QueryRow(ctx, getLatestWithdrawalSyncedBlockOfContractOnChain, arg.ContractAddress, arg.ChainID)
	var max interface{}
	err := row.Scan(&max)
	return max, err
}

const getWithdrawalByNullifierHash = `-- name: GetWithdrawalByNullifierHash :one
SELECT id, contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id FROM withdrawals WHERE nullifier_hash = $1
`

func (q *Queries) GetWithdrawalByNullifierHash(ctx context.Context, nullifierHash pgtype.Text) (Withdrawal, error) {
	row := q.db.QueryRow(ctx, getWithdrawalByNullifierHash, nullifierHash)
	var i Withdrawal
	err := row.Scan(
		&i.ID,
		&i.ContractAddress,
		&i.Recipient,
		&i.NullifierHash,
		&i.Relayer,
		&i.Fee,
		&i.TxHash,
		&i.Timestamp,
		&i.BlockNumber,
		&i.ChainID,
	)
	return i, err
}
