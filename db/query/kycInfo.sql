-- name: CreateKycInfo :one
INSERT INTO kyc_info (citizen_id, full_name, phone_number, date_of_birth, nationality, verifier, is_active, kyc_verified_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetKycInfoByCitizenID :one
SELECT * FROM kyc_info WHERE citizen_id = $1;

-- name: GetKycInfoByWalletAddress :one
SELECT k.* FROM kyc_info k
JOIN wallet_info w ON k.citizen_id = w.citizen_id
WHERE w.wallet_address = $1;

-- name: UpdateKycInfo :one
UPDATE kyc_info
SET full_name = $2, phone_number = $3, date_of_birth = $4, nationality = $5, verifier = $6, is_active = $7, kyc_verified_at = $8
WHERE citizen_id = $1
RETURNING *;