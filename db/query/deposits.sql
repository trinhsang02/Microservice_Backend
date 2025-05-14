
-- name: CreateDeposit :one
INSERT INTO deposits (contract_address, commitment, depositor, leaf_index, tx_hash, timestamp) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetAllDepositsOfContract :many
SELECT * FROM deposits WHERE contract_address = $1
ORDER BY leaf_index ASC;

-- name: getAllDepositsOfDepositor :many
SELECT * FROM deposits WHERE depositor = $1;
