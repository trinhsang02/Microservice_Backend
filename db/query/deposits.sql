
-- name: CreateDeposit :one
INSERT INTO deposits (contract_address, commitment, depositor, leaf_index, tx_hash, timestamp, block_number, chain_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetDepositsFromBlockToBlock :many
SELECT * FROM deposits 
WHERE contract_address = $1
AND chain_id = $2
AND block_number BETWEEN $3 AND $4
ORDER BY block_number ASC;

-- name: GetEarliestDepositSyncedBlock :one
SELECT MIN(block_number) FROM deposits 
WHERE contract_address = $1
AND chain_id = $2;

-- name: GetLatestDepositSyncedBlock :one
SELECT MAX(block_number) FROM deposits 
WHERE contract_address = $1
AND chain_id = $2;

-- name: GetDepositByCommitment :one
SELECT * FROM deposits 
WHERE commitment = $1;

-- name: GetLeaves :many
SELECT commitment FROM deposits
WHERE contract_address = $1
AND chain_id = $2
ORDER BY leaf_index ASC;