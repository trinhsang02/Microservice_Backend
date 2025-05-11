package database

import (
	"database/sql"
	"fmt"
)

// ----- DEPOSIT OPERATIONS -----

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

// UpdateDepositLeafIndex cập nhật leaf index cho deposit
func UpdateDepositLeafIndex(id int, leafIndex int) error {
	query := `
		UPDATE DEPOSITS
		SET leaf_index = $2
		WHERE id = $1
	`

	_, err := DB.Exec(query, id, leafIndex)
	if err != nil {
		return fmt.Errorf("không thể cập nhật leaf index: %v", err)
	}

	return nil
}

// ----- WITHDRAWAL OPERATIONS -----

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

// UpdateWithdrawalRelayer cập nhật thông tin relayer cho withdrawal
func UpdateWithdrawalRelayer(id int, relayer string, fee float64) error {
	query := `
		UPDATE WITHDRAWALS
		SET relayer = $2, fee = $3
		WHERE id = $1
	`

	_, err := DB.Exec(query, id, relayer, fee)
	if err != nil {
		return fmt.Errorf("không thể cập nhật thông tin relayer: %v", err)
	}

	return nil
}

// ----- KYC OPERATIONS -----

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

// IsKYCVerified kiểm tra xem một địa chỉ ví đã được xác minh KYC chưa
func IsKYCVerified(walletAddress string) (bool, error) {
	kyc, err := GetKYCByWalletAddress(walletAddress)
	if err != nil {
		return false, err
	}

	if kyc == nil {
		return false, nil
	}

	return kyc.IsActive, nil
}
