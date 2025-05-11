CREATE INDEX IF NOT EXISTS idx_deposits_commitment ON deposits(commitment);
CREATE INDEX IF NOT EXISTS idx_deposits_depositor ON deposits(depositor);
CREATE INDEX IF NOT EXISTS idx_withdrawals_nullifier_hash ON withdrawals(nullifier_hash);
CREATE INDEX IF NOT EXISTS idx_kyc_wallet_address ON kyc(wallet_address);

ALTER TABLE deposits OWNER TO postgres;
ALTER TABLE withdrawals OWNER TO postgres;
ALTER TABLE kyc OWNER TO postgres;