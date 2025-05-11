package database

import (
	"database/sql"
	"fmt"
)

// ----- DEPOSIT OPERATIONS -----

// GetDepositByID lấy thông tin deposit theo ID
func GetDepositByID(id int) (*Deposit, error) {
	query := `
		SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash
		FROM DEPOSITS
		WHERE id = $1
	`

	var deposit Deposit
	err := DB.QueryRow(query, id).Scan(
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
		return nil, fmt.Errorf("không thể tìm deposit theo ID: %v", err)
	}

	return &deposit, nil
}

// GetDepositsByDepositor lấy danh sách deposits theo địa chỉ người gửi
func GetDepositsByDepositor(depositor string, limit, offset int) ([]Deposit, error) {
	query := `
		SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash
		FROM DEPOSITS
		WHERE depositor = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.Query(query, depositor, limit, offset)
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

// ----- WITHDRAWAL OPERATIONS -----

// GetWithdrawalByID lấy thông tin withdrawal theo ID
func GetWithdrawalByID(id int) (*Withdrawal, error) {
	query := `
		SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash
		FROM WITHDRAWALS
		WHERE id = $1
	`

	var withdrawal Withdrawal
	err := DB.QueryRow(query, id).Scan(
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
		return nil, fmt.Errorf("không thể tìm withdrawal theo ID: %v", err)
	}

	return &withdrawal, nil
}

// GetWithdrawalsByRecipient lấy danh sách withdrawals theo địa chỉ người nhận
func GetWithdrawalsByRecipient(recipient string, limit, offset int) ([]Withdrawal, error) {
	query := `
		SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash
		FROM WITHDRAWALS
		WHERE recipient = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.Query(query, recipient, limit, offset)
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

// ----- KYC OPERATIONS -----

// GetKYCByWalletAddress lấy thông tin KYC theo địa chỉ ví
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
