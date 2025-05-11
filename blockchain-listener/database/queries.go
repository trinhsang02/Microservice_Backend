package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type BlockchainEvent struct {
	ID              int             `json:"id"`
	EventType       string          `json:"event_type"`
	TransactionHash string          `json:"transaction_hash"`
	BlockNumber     int64           `json:"block_number"`
	EventData       json.RawMessage `json:"event_data"`
	CreatedAt       time.Time       `json:"created_at"`
}


type ProcessedBlock struct {
	ID          int       `json:"id"`
	BlockNumber int64     `json:"block_number"`
	BlockHash   string    `json:"block_hash"`
	ProcessedAt time.Time `json:"processed_at"`
}

func SaveBlockchainEvent(eventType, txHash string, blockNumber int64, eventData json.RawMessage) (int, error) {
	var id int
	query := `
		INSERT INTO blockchain_events (event_type, transaction_hash, block_number, event_data)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := DB.QueryRow(query, eventType, txHash, blockNumber, eventData).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("không thể lưu blockchain event: %v", err)
	}
	return id, nil
}

// GetBlockchainEvents lấy danh sách blockchain events từ database
func GetBlockchainEvents(limit, offset int) ([]BlockchainEvent, error) {
	query := `
		SELECT id, event_type, transaction_hash, block_number, event_data, created_at
		FROM blockchain_events
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy blockchain events: %v", err)
	}
	defer rows.Close()

	var events []BlockchainEvent
	for rows.Next() {
		var event BlockchainEvent
		err := rows.Scan(
			&event.ID,
			&event.EventType,
			&event.TransactionHash,
			&event.BlockNumber,
			&event.EventData,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể xử lý blockchain event: %v", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc blockchain events: %v", err)
	}

	return events, nil
}

// GetEventsByType lấy events theo loại
func GetEventsByType(eventType string, limit, offset int) ([]BlockchainEvent, error) {
	query := `
		SELECT id, event_type, transaction_hash, block_number, event_data, created_at
		FROM blockchain_events
		WHERE event_type = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.Query(query, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy blockchain events theo loại: %v", err)
	}
	defer rows.Close()

	var events []BlockchainEvent
	for rows.Next() {
		var event BlockchainEvent
		err := rows.Scan(
			&event.ID,
			&event.EventType,
			&event.TransactionHash,
			&event.BlockNumber,
			&event.EventData,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("không thể xử lý blockchain event: %v", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình đọc blockchain events: %v", err)
	}

	return events, nil
}

// SaveProcessedBlock lưu thông tin của block đã xử lý
func SaveProcessedBlock(blockNumber int64, blockHash string) (int, error) {
	var id int
	query := `
		INSERT INTO processed_blocks (block_number, block_hash)
		VALUES ($1, $2)
		ON CONFLICT (block_number) DO UPDATE
		SET block_hash = $2, processed_at = CURRENT_TIMESTAMP
		RETURNING id
	`
	err := DB.QueryRow(query, blockNumber, blockHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("không thể lưu processed block: %v", err)
	}
	return id, nil
}

// IsBlockProcessed kiểm tra xem một block đã được xử lý chưa
func IsBlockProcessed(blockNumber int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM processed_blocks WHERE block_number = $1)
	`
	err := DB.QueryRow(query, blockNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("không thể kiểm tra block đã xử lý: %v", err)
	}
	return exists, nil
}

// GetLastProcessedBlock lấy block được xử lý gần đây nhất
func GetLastProcessedBlock() (*ProcessedBlock, error) {
	query := `
		SELECT id, block_number, block_hash, processed_at
		FROM processed_blocks
		ORDER BY block_number DESC
		LIMIT 1
	`

	var block ProcessedBlock
	err := DB.QueryRow(query).Scan(
		&block.ID,
		&block.BlockNumber,
		&block.BlockHash,
		&block.ProcessedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("không thể lấy block được xử lý gần đây nhất: %v", err)
	}

	return &block, nil
}
