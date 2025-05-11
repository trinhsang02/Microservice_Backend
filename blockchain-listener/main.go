package main

import (
	"blockchain-listener/database"
	"blockchain-listener/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// BlockchainEvent cấu trúc dữ liệu của event để gửi qua RabbitMQ
type BlockchainEvent struct {
	EventType       string          `json:"event_type"`
	TransactionHash string          `json:"transaction_hash"`
	BlockNumber     int64           `json:"block_number"`
	EventData       json.RawMessage `json:"event_data"`
}

func main() {
	// Khởi tạo kết nối PostgreSQL
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Kiểm tra kết nối bằng cách lấy danh sách deposits
	deposits, err := database.GetDeposits(5, 0)
	if err != nil {
		log.Printf("Warning: Failed to get deposits: %v", err)
	} else {
		log.Printf("Successfully retrieved %d deposits", len(deposits))
		for _, deposit := range deposits {
			log.Printf("Deposit ID: %d, Commitment: %s, Depositor: %s",
				deposit.ID, deposit.Commitment, deposit.Depositor)
		}
	}

	// Kiểm tra kết nối bằng cách lấy danh sách withdrawals
	withdrawals, err := database.GetWithdrawals(5, 0)
	if err != nil {
		log.Printf("Warning: Failed to get withdrawals: %v", err)
	} else {
		log.Printf("Successfully retrieved %d withdrawals", len(withdrawals))
		for _, withdrawal := range withdrawals {
			log.Printf("Withdrawal ID: %d, Recipient: %s, Nullifier Hash: %s",
				withdrawal.ID, withdrawal.Recipient, withdrawal.NullifierHash)
		}
	}

	// Kiểm tra kết nối bằng cách lấy danh sách KYC
	kycRecords, err := database.GetAllKYC(5, 0)
	if err != nil {
		log.Printf("Warning: Failed to get KYC records: %v", err)
	} else {
		log.Printf("Successfully retrieved %d KYC records", len(kycRecords))
		for _, kyc := range kycRecords {
			log.Printf("KYC CitizenID: %s, Full Name: %s, Wallet Address: %s",
				kyc.CitizenID, kyc.FullName, kyc.WalletAddress)
		}
	}

	// RabbitMQ connection details
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	queueName := "blockchain_events"

	// Initialize the producer
	producer, err := rabbitmq.NewProducer(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	// Publish messages periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Xử lý khi chương trình bị kết thúc
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	log.Println("Blockchain Listener is publishing messages...")
	// 	for t := range ticker.C {
	// 		// Mô phỏng việc phát hiện một giao dịch deposit trên blockchain
	// 		handleDepositEvent(producer, t)

	// 		// Cứ 3 lần deposit thì tạo 1 withdrawal để mô phỏng
	// 		if t.Second()%15 < 5 {
	// 			handleWithdrawalEvent(producer, t)
	// 		}
	// 	}
	// }()

	// Chờ tín hiệu kết thúc
	<-sigChan
	log.Println("Shutting down gracefully...")
}

// handleDepositEvent xử lý và lưu trữ một sự kiện deposit
// func handleDepositEvent(producer *rabbitmq.Producer, t time.Time) {
// 	// Tạo dữ liệu mô phỏng cho một deposit transaction
// 	commitment := "commit_" + time.Now().Format("20060102150405")
// 	depositor := "0x" + fmt.Sprintf("%x", t.Nanosecond())[0:8]
// 	leafIndex := int(t.Unix() % 1000)
// 	txHash := "0xd_" + time.Now().Format("20060102150405")

// 	// Lưu vào database
// 	id, err := database.CreateDeposit(commitment, depositor, leafIndex, txHash)
// 	if err != nil {
// 		log.Printf("Failed to create deposit: %v", err)
// 		return
// 	}
// 	log.Printf("Created deposit with ID: %d", id)

// 	// Tạo dữ liệu JSON cho event
// 	eventData := map[string]interface{}{
// 		"commitment": commitment,
// 		"depositor":  depositor,
// 		"leaf_index": leafIndex,
// 		"amount":     "1.5",
// 		"token":      "ETH",
// 		"timestamp":  time.Now().Format(time.RFC3339),
// 	}
// 	eventDataBytes, _ := json.Marshal(eventData)

// 	// Tạo thông điệp để gửi qua RabbitMQ
// 	event := BlockchainEvent{
// 		EventType:       "Deposit",
// 		TransactionHash: txHash,
// 		BlockNumber:     int64(1000000 + t.Unix()%10000),
// 		EventData:       eventDataBytes,
// 	}

// 	// Chuyển event thành JSON
// 	message, err := json.Marshal(event)
// 	if err != nil {
// 		log.Printf("Failed to marshal deposit event: %v", err)
// 		return
// 	}

// 	// Gửi message qua RabbitMQ
// 	err = producer.Publish(string(message))
// 	if err != nil {
// 		log.Printf("Failed to publish deposit event: %v", err)
// 		return
// 	}

// 	log.Printf("Published deposit event with commitment: %s", commitment)
// }

// handleWithdrawalEvent xử lý và lưu trữ một sự kiện withdrawal
func handleWithdrawalEvent(producer *rabbitmq.Producer, t time.Time) {
	// Tạo dữ liệu mô phỏng cho một withdrawal transaction
	recipient := "0x" + fmt.Sprintf("%x", t.Nanosecond())[0:8]
	nullifierHash := "null_" + time.Now().Format("20060102150405")
	relayer := "0xrelayer_" + fmt.Sprintf("%x", t.Unix())[0:4]
	fee := float64(t.Unix()%100) / 1000.0
	txHash := "0xw_" + time.Now().Format("20060102150405")

	// Lưu vào database
	id, err := database.CreateWithdrawal(recipient, nullifierHash, relayer, fee, txHash)
	if err != nil {
		log.Printf("Failed to create withdrawal: %v", err)
		return
	}
	log.Printf("Created withdrawal with ID: %d", id)

	// Tạo dữ liệu JSON cho event
	eventData := map[string]interface{}{
		"recipient":      recipient,
		"nullifier_hash": nullifierHash,
		"relayer":        relayer,
		"fee":            fee,
		"amount":         "1.0",
		"token":          "ETH",
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	eventDataBytes, _ := json.Marshal(eventData)

	// Tạo thông điệp để gửi qua RabbitMQ
	event := BlockchainEvent{
		EventType:       "Withdrawal",
		TransactionHash: txHash,
		BlockNumber:     int64(1000000 + t.Unix()%10000),
		EventData:       eventDataBytes,
	}

	// Chuyển event thành JSON
	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal withdrawal event: %v", err)
		return
	}

	// Gửi message qua RabbitMQ
	err = producer.Publish(string(message))
	if err != nil {
		log.Printf("Failed to publish withdrawal event: %v", err)
		return
	}

	log.Printf("Published withdrawal event for recipient: %s", recipient)
}
