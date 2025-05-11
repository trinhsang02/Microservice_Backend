-- Kiểm tra và tạo database nếu chưa tồn tại (cú pháp PostgreSQL)
SELECT 'CREATE DATABASE microservices'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'microservices');

-- Kết nối đến database microservices
\c microservices

-- Đảm bảo schema public được sử dụng
SET search_path TO public;

-- Tạo bảng với tên viết thường
CREATE TABLE IF NOT EXISTS deposits (
    id SERIAL PRIMARY KEY,
    commitment VARCHAR,
    depositor VARCHAR,
    leaf_index INT,
    timestamp TIMESTAMP,
    tx_hash VARCHAR
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    recipient VARCHAR,
    nullifier_hash VARCHAR,
    relayer VARCHAR,
    fee NUMERIC,
    timestamp TIMESTAMP,
    tx_hash VARCHAR
);

CREATE TABLE IF NOT EXISTS kyc (
    citizen_id VARCHAR PRIMARY KEY,
    wallet_address VARCHAR,
    full_name VARCHAR,
    phone_number VARCHAR,
    date_of_birth DATE,
    nationality VARCHAR,
    kyc_verified_at TIMESTAMP,
    verifier VARCHAR,
    is_active BOOLEAN,
    wallet_signature VARCHAR
);