-- Đảm bảo kết nối đến database microservices
\c microservices

-- Thêm dữ liệu mẫu vào bảng deposits
INSERT INTO deposits (commitment, depositor, leaf_index, timestamp, tx_hash)
VALUES 
('0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', '0xA1B2C3D4E5F67890', 0, NOW() - INTERVAL '3 day', '0xabcd1234'),
('0x2345678901abcdef2345678901abcdef2345678901abcdef2345678901abcdef', '0xB2C3D4E5F6789012', 1, NOW() - INTERVAL '2 day', '0xbcde2345'),
('0x3456789012abcdef3456789012abcdef3456789012abcdef3456789012abcdef', '0xC3D4E5F67890123A', 2, NOW() - INTERVAL '1 day', '0xcdef3456');

-- Thêm dữ liệu mẫu vào bảng withdrawals
INSERT INTO withdrawals (recipient, nullifier_hash, relayer, fee, timestamp, tx_hash)
VALUES 
('0xD4E5F6789012345B', '0x4567890123abcdef4567890123abcdef4567890123abcdef4567890123abcdef', '0xRelay123456789', 0.01, NOW() - INTERVAL '3 day', '0xdef45678'),
('0xE5F6789012345C6', '0x5678901234abcdef5678901234abcdef5678901234abcdef5678901234abcdef', '0xRelay234567890', 0.02, NOW() - INTERVAL '2 day', '0xefab5678'),
('0xF6789012345D67E', '0x6789012345abcdef6789012345abcdef6789012345abcdef6789012345abcdef', '0xRelay345678901', 0.03, NOW() - INTERVAL '1 day', '0xfabc6789');

-- Thêm dữ liệu mẫu vào bảng kyc
INSERT INTO kyc (citizen_id, wallet_address, full_name, phone_number, date_of_birth, nationality, kyc_verified_at, verifier, is_active, wallet_signature)
VALUES 
('123456789012', '0xA1B2C3D4E5F67890', 'Nguyen Van A', '0901234567', '1990-01-01', 'Vietnamese', NOW() - INTERVAL '10 day', '0xVerifier1', true, '0xsig1234'),
('234567890123', '0xB2C3D4E5F6789012', 'Tran Thi B', '0912345678', '1992-02-02', 'Vietnamese', NOW() - INTERVAL '8 day', '0xVerifier2', true, '0xsig2345'),
('345678901234', '0xC3D4E5F67890123A', 'Le Van C', '0923456789', '1995-03-03', 'Vietnamese', NOW() - INTERVAL '5 day', '0xVerifier3', false, '0xsig3456');