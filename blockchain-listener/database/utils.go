package database

import (
	"fmt"
	"log"
	"time"
)

// DBUtils chứa các hàm tiện ích để thao tác với database
type DBUtils struct{}

// NewDBUtils tạo một instance mới của DBUtils
func NewDBUtils() *DBUtils {
	return &DBUtils{}
}

// GetRecentDeposits trả về danh sách N giao dịch deposit mới nhất
func (u *DBUtils) GetRecentDeposits(limit int) ([]Deposit, error) {
	query := `SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash 
			  FROM deposits 
			  ORDER BY timestamp DESC 
			  LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn deposits: %v", err)
	}
	defer rows.Close()

	var deposits []Deposit
	for rows.Next() {
		var d Deposit
		if err := rows.Scan(&d.ID, &d.Commitment, &d.Depositor, &d.LeafIndex, &d.Timestamp, &d.TxHash); err != nil {
			return nil, fmt.Errorf("lỗi khi scan deposit: %v", err)
		}
		deposits = append(deposits, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi sau khi truy vấn deposits: %v", err)
	}

	return deposits, nil
}

// GetRecentWithdrawals trả về danh sách N giao dịch withdrawal mới nhất
func (u *DBUtils) GetRecentWithdrawals(limit int) ([]Withdrawal, error) {
	query := `SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash 
			  FROM withdrawals 
			  ORDER BY timestamp DESC 
			  LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn withdrawals: %v", err)
	}
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var w Withdrawal
		if err := rows.Scan(&w.ID, &w.Recipient, &w.NullifierHash, &w.Relayer, &w.Fee, &w.Timestamp, &w.TxHash); err != nil {
			return nil, fmt.Errorf("lỗi khi scan withdrawal: %v", err)
		}
		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi sau khi truy vấn withdrawals: %v", err)
	}

	return withdrawals, nil
}

// GetUserKYCStatus kiểm tra tình trạng KYC của một địa chỉ ví
func (u *DBUtils) GetUserKYCStatus(walletAddress string) (*KYC, error) {
	query := `SELECT citizen_id, wallet_address, full_name, phone_number, date_of_birth, 
			  nationality, kyc_verified_at, verifier, is_active, wallet_signature 
			  FROM kyc 
			  WHERE wallet_address = $1`

	var kyc KYC
	err := DB.QueryRow(query, walletAddress).Scan(
		&kyc.CitizenID, &kyc.WalletAddress, &kyc.FullName, &kyc.PhoneNumber, &kyc.DateOfBirth,
		&kyc.Nationality, &kyc.KYCVerifiedAt, &kyc.Verifier, &kyc.IsActive, &kyc.WalletSignature,
	)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn thông tin KYC: %v", err)
	}

	return &kyc, nil
}

// GetDepositsByWalletAddress trả về các giao dịch deposit của một địa chỉ ví
func (u *DBUtils) GetDepositsByWalletAddress(walletAddress string) ([]Deposit, error) {
	query := `SELECT id, commitment, depositor, leaf_index, timestamp, tx_hash 
			  FROM deposits 
			  WHERE depositor = $1 
			  ORDER BY timestamp DESC`

	rows, err := DB.Query(query, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn deposits theo địa chỉ ví: %v", err)
	}
	defer rows.Close()

	var deposits []Deposit
	for rows.Next() {
		var d Deposit
		if err := rows.Scan(&d.ID, &d.Commitment, &d.Depositor, &d.LeafIndex, &d.Timestamp, &d.TxHash); err != nil {
			return nil, fmt.Errorf("lỗi khi scan deposit: %v", err)
		}
		deposits = append(deposits, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi sau khi truy vấn deposits: %v", err)
	}

	return deposits, nil
}

// GetWithdrawalsByWalletAddress trả về các giao dịch withdrawal của một địa chỉ ví
func (u *DBUtils) GetWithdrawalsByWalletAddress(walletAddress string) ([]Withdrawal, error) {
	query := `SELECT id, recipient, nullifier_hash, relayer, fee, timestamp, tx_hash 
			  FROM withdrawals 
			  WHERE recipient = $1 
			  ORDER BY timestamp DESC`

	rows, err := DB.Query(query, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn withdrawals theo địa chỉ ví: %v", err)
	}
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var w Withdrawal
		if err := rows.Scan(&w.ID, &w.Recipient, &w.NullifierHash, &w.Relayer, &w.Fee, &w.Timestamp, &w.TxHash); err != nil {
			return nil, fmt.Errorf("lỗi khi scan withdrawal: %v", err)
		}
		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi sau khi truy vấn withdrawals: %v", err)
	}

	return withdrawals, nil
}

// GetTransactionStats lấy thống kê số lượng giao dịch trong một khoảng thời gian
func (u *DBUtils) GetTransactionStats(startTime, endTime time.Time) (map[string]int, error) {
	stats := make(map[string]int)

	// Đếm số lượng deposits
	depositQuery := `SELECT COUNT(*) FROM deposits WHERE timestamp BETWEEN $1 AND $2`
	var depositCount int
	if err := DB.QueryRow(depositQuery, startTime, endTime).Scan(&depositCount); err != nil {
		return nil, fmt.Errorf("lỗi khi đếm deposits: %v", err)
	}
	stats["deposits"] = depositCount

	// Đếm số lượng withdrawals
	withdrawalQuery := `SELECT COUNT(*) FROM withdrawals WHERE timestamp BETWEEN $1 AND $2`
	var withdrawalCount int
	if err := DB.QueryRow(withdrawalQuery, startTime, endTime).Scan(&withdrawalCount); err != nil {
		return nil, fmt.Errorf("lỗi khi đếm withdrawals: %v", err)
	}
	stats["withdrawals"] = withdrawalCount

	return stats, nil
}

// PrintDatabaseStats in thống kê cơ bản về dữ liệu trong database
func (u *DBUtils) PrintDatabaseStats() {
	var depositCount, withdrawalCount, kycCount int

	// Đếm số lượng deposits
	if err := DB.QueryRow("SELECT COUNT(*) FROM deposits").Scan(&depositCount); err != nil {
		log.Printf("Lỗi khi đếm deposits: %v", err)
	}

	// Đếm số lượng withdrawals
	if err := DB.QueryRow("SELECT COUNT(*) FROM withdrawals").Scan(&withdrawalCount); err != nil {
		log.Printf("Lỗi khi đếm withdrawals: %v", err)
	}

	// Đếm số lượng người dùng đã KYC
	if err := DB.QueryRow("SELECT COUNT(*) FROM kyc").Scan(&kycCount); err != nil {
		log.Printf("Lỗi khi đếm KYC: %v", err)
	}

	fmt.Println("=== THỐNG KÊ CƠ SỞ DỮ LIỆU ===")
	fmt.Printf("Tổng số deposits: %d\n", depositCount)
	fmt.Printf("Tổng số withdrawals: %d\n", withdrawalCount)
	fmt.Printf("Tổng số người dùng đã KYC: %d\n", kycCount)
	fmt.Println("==============================")
}

// TestDatabaseConnection kiểm tra kết nối đến database
func (u *DBUtils) TestDatabaseConnection() error {
	// Thực hiện một truy vấn đơn giản để kiểm tra kết nối
	var result int
	err := DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("lỗi kết nối database: %v", err)
	}

	if result == 1 {
		fmt.Println("✅ Kết nối database thành công!")
		return nil
	}

	return fmt.Errorf("kiểm tra kết nối database thất bại")
}
