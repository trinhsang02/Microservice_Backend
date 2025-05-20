CREATE TABLE IF NOT EXISTS deposits (
    id SERIAL PRIMARY KEY,
    contract_address VARCHAR(255),
    commitment VARCHAR(255),
    depositor VARCHAR(255),
    leaf_index INT,
    tx_hash VARCHAR(255),
    timestamp NUMERIC,
    block_number INT,
    chain_id INT,
    UNIQUE (commitment)
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    contract_address VARCHAR(255),
    recipient VARCHAR(255),
    nullifier_hash VARCHAR(255),
    relayer VARCHAR(255),
    fee NUMERIC,
    tx_hash VARCHAR(255),
    timestamp NUMERIC,
    block_number INT,
    chain_id INT,
    UNIQUE (nullifier_hash)
);

CREATE TABLE IF NOT EXISTS kyc_info (
    citizen_id VARCHAR(255) PRIMARY KEY,
    full_name VARCHAR(255),
    phone_number VARCHAR(50),
    date_of_birth DATE,
    nationality VARCHAR(100),
    verifier VARCHAR(255),
    is_active BOOLEAN,
    kyc_verified_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS wallet_info (
    wallet_address VARCHAR(255) PRIMARY KEY,
    citizen_id VARCHAR(255) REFERENCES kyc_info(citizen_id),
    wallet_signature VARCHAR(255), -- sign on citizen_id
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE INDEX IF NOT EXISTS contract_address_chain_id_index ON deposits (contract_address, chain_id);
CREATE INDEX IF NOT EXISTS citizen_id_index ON kyc_info (citizen_id);

