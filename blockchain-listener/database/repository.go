package database

import (
	"database/sql"
	"fmt"
	"time"
)

// ----- DEPOSIT OPERATIONS -----

// CreateDeposit thêm một giao dịch deposit mới vào database
func CreateDeposit(commitment, depositor string, leafIndex int, txHash string) (int, error) {
	var id int
	query := `
		INSERT INTO DEPOSITS (commitment, depositor, leaf_index, timestamp, tx_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := DB.QueryRow(query, commitment, depositor, leafIndex, time.Now(), txHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("không thể tạo deposit: %v", err)
	}
	return id, nil
}

// GetDeposits lấy danh sách deposits từ database
func GetDeposits(limit, offset int) ([]Deposit, error) {
	query := `
		SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash
		FROM DEPOSITS
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách deposits: %v", err)
	}
	defer rows.Close()
	
	var deposits []Deposit
	for rows.Next() {
		var deposit Deposit
		err := rows.Scan(
			&deposit.ID,
			&deposit.Commitment,
			&deposit.Depositor,
			&deposit.LeafIndex,
			&deposit.Timestamp,
			&deposit.TxHash,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu deposit: %v", err)
		}
		deposits = append(deposits, deposit)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc deposits: %v", err)
	}
	
	return deposits, nil
}

// GetDepositByCommitment tìm deposit theo commitment
func GetDepositByCommitment(commitment string) (*Deposit, error) {
	query := `
		SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash
		FROM DEPOSITS
		WHERE commitment = $1
	`
	
	var deposit Deposit
	err := DB.QueryRow(query, commitment).Scan(
		&deposit.ID,
		&deposit.Commitment,
		&deposit.Depositor,
		&deposit.LeafIndex,
		&deposit.Timestamp,
		&deposit.TxHash,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Không tìm thấy
	}
	
	if err != nil {
		return nil, fmt.Errorf("không thể tìm deposit theo commitment: %v", err)
	}
	
	return &deposit, nil
}

// GetDepositByDepositor tìm deposits theo địa chỉ depositor
func GetDepositsByDepositor(depositor string, limit, offset int) ([]Deposit, error) {
	query := `
		SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash
		FROM DEPOSITS
		WHERE depositor = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := DB.Query(query, depositor, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy deposits theo depositor: %v", err)
	}
	defer rows.Close()
	
	var deposits []Deposit
	for rows.Next() {
		var deposit Deposit
		err := rows.Scan(
			&deposit.ID,
			&deposit.Commitment,
			&deposit.Depositor,
			&deposit.LeafIndex,
			&deposit.Timestamp,
			&deposit.TxHash,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu deposit: %v", err)
		}
		deposits = append(deposits, deposit)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc deposits: %v", err)
	}
	
	return deposits, nil
}

// ----- WITHDRAWAL OPERATIONS -----

// CreateWithdrawal thêm một giao dịch withdrawal mới vào database
func CreateWithdrawal(recipient, nullifierHash, relayer string, fee float64, txHash string) (int, error) {
	var id int
	query := `
		INSERT INTO WITHDRAWALS (recipient, nullifier_hash, relayer, fee, timestamp, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := DB.QueryRow(query, recipient, nullifierHash, relayer, fee, time.Now(), txHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("không thể tạo withdrawal: %v", err)
	}
	return id, nil
}

// GetWithdrawals lấy danh sách withdrawals từ database
func GetWithdrawals(limit, offset int) ([]Withdrawal, error) {
	query := `
		SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash
		FROM WITHDRAWALS
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách withdrawals: %v", err)
	}
	defer rows.Close()
	
	var withdrawals []Withdrawal
	for rows.Next() {
		var withdrawal Withdrawal
		err := rows.Scan(
			&withdrawal.ID,
			&withdrawal.Recipient,
			&withdrawal.NullifierHash,
			&withdrawal.Relayer,
			&withdrawal.Fee,
			&withdrawal.Timestamp,
			&withdrawal.TxHash,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu withdrawal: %v", err)
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc withdrawals: %v", err)
	}
	
	return withdrawals, nil
}

// GetWithdrawalByNullifierHash tìm withdrawal theo nullifier hash
func GetWithdrawalByNullifierHash(nullifierHash string) (*Withdrawal, error) {
	query := `
		SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash
		FROM WITHDRAWALS
		WHERE nullifier_hash = $1
	`
	
	var withdrawal Withdrawal
	err := DB.QueryRow(query, nullifierHash).Scan(
		&withdrawal.ID,
		&withdrawal.Recipient,
		&withdrawal.NullifierHash,
		&withdrawal.Relayer,
		&withdrawal.Fee,
		&withdrawal.Timestamp,
		&withdrawal.TxHash,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Không tìm thấy
	}
	
	if err != nil {
		return nil, fmt.Errorf("không thể tìm withdrawal theo nullifier hash: %v", err)
	}
	
	return &withdrawal, nil
}

// GetWithdrawalsByRecipient tìm withdrawals theo recipient
func GetWithdrawalsByRecipient(recipient string, limit, offset int) ([]Withdrawal, error) {
	query := `
		SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash
		FROM WITHDRAWALS
		WHERE recipient = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := DB.Query(query, recipient, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy withdrawals theo recipient: %v", err)
	}
	defer rows.Close()
	
	var withdrawals []Withdrawal
	for rows.Next() {
		var withdrawal Withdrawal
		err := rows.Scan(
			&withdrawal.ID,
			&withdrawal.Recipient,
			&withdrawal.NullifierHash,
			&withdrawal.Relayer,
			&withdrawal.Fee,
			&withdrawal.Timestamp,
			&withdrawal.TxHash,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu withdrawal: %v", err)
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc withdrawals: %v", err)
	}
	
	return withdrawals, nil
}

// ----- KYC OPERATIONS -----

// CreateKYC thêm một bản ghi KYC mới vào database
func CreateKYC(citizenID, walletAddress, fullName, phoneNumber string, dateOfBirth time.Time, 
    nationality, verifier, walletSignature string, isActive bool) error {
	query := `
		INSERT INTO KYC (citizen_id, wallet_address, full_name, phone_number, date_of_birth, 
		    nationality, kyc_verified_at, verifier, is_active, wallet_signature)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := DB.Exec(query, citizenID, walletAddress, fullName, phoneNumber, dateOfBirth, 
	    nationality, time.Now(), verifier, isActive, walletSignature)
	if err != nil {
		return fmt.Errorf("không thể tạo KYC record: %v", err)
	}
	return nil
}

// GetKYCByCitizenID lấy thông tin KYC theo citizenID
func GetKYCByCitizenID(citizenID string) (*KYC, error) {
	query := `
		SELECT citizen_id, wallet_address, full_name, phone_number, date_of_birth, 
		    nationality, kyc_verified_at, verifier, is_active, wallet_signature
		FROM KYC
		WHERE citizen_id = $1
	`
	
	var kyc KYC
	err := DB.QueryRow(query, citizenID).Scan(
		&kyc.CitizenID,
		&kyc.WalletAddress,
		&kyc.FullName,
		&kyc.PhoneNumber,
		&kyc.DateOfBirth,
		&kyc.Nationality,
		&kyc.KYCVerifiedAt,
		&kyc.Verifier,
		&kyc.IsActive,
		&kyc.WalletSignature,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Không tìm thấy
	}
	
	if err != nil {
		return nil, fmt.Errorf("không thể tìm KYC theo citizenID: %v", err)
	}
	
	return &kyc, nil
}

// GetKYCByWalletAddress lấy thông tin KYC theo wallet address
func GetKYCByWalletAddress(walletAddress string) (*KYC, error) {
	query := `
		SELECT citizen_id, wallet_address, full_name, phone_number, date_of_birth, 
		    nationality, kyc_verified_at, verifier, is_active, wallet_signature
		FROM KYC
		WHERE wallet_address = $1
	`
	
	var kyc KYC
	err := DB.QueryRow(query, walletAddress).Scan(
		&kyc.CitizenID,
		&kyc.WalletAddress,
		&kyc.FullName,
		&kyc.PhoneNumber,
		&kyc.DateOfBirth,
		&kyc.Nationality,
		&kyc.KYCVerifiedAt,
		&kyc.Verifier,
		&kyc.IsActive,
		&kyc.WalletSignature,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Không tìm thấy
	}
	
	if err != nil {
		return nil, fmt.Errorf("không thể tìm KYC theo wallet address: %v", err)
	}
	
	return &kyc, nil
}

// UpdateKYCStatus cập nhật trạng thái KYC
func UpdateKYCStatus(citizenID string, isActive bool) error {
	query := `
		UPDATE KYC
		SET is_active = $2
		WHERE citizen_id = $1
	`
	
	_, err := DB.Exec(query, citizenID, isActive)
	if err != nil {
		return fmt.Errorf("không thể cập nhật trạng thái KYC: %v", err)
	}
	
	return nil
}

// GetAllKYC lấy danh sách tất cả KYC
func GetAllKYC(limit, offset int) ([]KYC, error) {
	query := `
		SELECT citizen_id, wallet_address, full_name, phone_number, date_of_birth, 
		    nationality, kyc_verified_at, verifier, is_active, wallet_signature
		FROM KYC
		ORDER BY kyc_verified_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách KYC: %v", err)
	}
	defer rows.Close()
	
	var kycList []KYC
	for rows.Next() {
		var kyc KYC
		err := rows.Scan(
			&kyc.CitizenID,
			&kyc.WalletAddress,
			&kyc.FullName,
			&kyc.PhoneNumber,
			&kyc.DateOfBirth,
			&kyc.Nationality,
			&kyc.KYCVerifiedAt,
			&kyc.Verifier,
			&kyc.IsActive,
			&kyc.WalletSignature,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu KYC: %v", err)
		}
		kycList = append(kycList, kyc)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc KYC: %v", err)
	}
	
	return kycList, nil
}