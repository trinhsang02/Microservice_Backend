
-- name: CreateWithdrawal :one
INSERT INTO withdrawals (contract_address, recipient, nullifier_hash, relayer, fee, tx_hash, timestamp, block_number, chain_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetAllWithdrawalsOfContract :many
SELECT * FROM withdrawals WHERE contract_address = $1
ORDER BY timestamp ASC;

-- name: GetAllWithdrawalsOfRecipient :many
SELECT * FROM withdrawals WHERE recipient = $1
ORDER BY timestamp ASC;

-- name: GetLatestWithdrawalSyncedBlockOfContractOnChain :one
SELECT MAX(block_number) FROM withdrawals WHERE contract_address = $1 AND chain_id = $2;

-- name: GetWithdrawalByNullifierHash :one
SELECT * FROM withdrawals WHERE nullifier_hash = $1;