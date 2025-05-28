-- name: CreateWalletInfo :one
INSERT INTO wallet_info (
    wallet_address,
    citizen_id,
    wallet_signature,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: CreateOrUpdateWalletInfo :exec
INSERT INTO wallet_info (wallet_address, citizen_id, wallet_signature, created_at)
VALUES ($1, $2, $3, now())
ON CONFLICT (wallet_address) DO UPDATE
SET citizen_id = EXCLUDED.citizen_id, 
    wallet_signature = EXCLUDED.wallet_signature,
    created_at = now();