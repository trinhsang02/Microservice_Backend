package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourusername/yourrepo/db/sqlc"
	"github.com/yourusername/yourrepo/mq/rabbitmq"
)

type MintMessage struct {
	sqlc.KycInfo  `json:"kyc"`
	WalletAddress string `json:"wallet_address"`
}

const kycWalletNFTABI = `[{"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"mint","stateMutability":"nonpayable","type":"function"}]`

func StartMintWorker(repo *sqlc.Repository, rabbitmqURL string) {
	log.Println("Starting mint worker...")
	consumer, err := rabbitmq.NewConsumer(rabbitmqURL, "kyc-mint-exchange", "topic", "kyc-mint-queue", []string{"kyc.mint"})
	if err != nil {
		log.Printf("Failed to create consumer: %v", err)
		return
	}
	defer consumer.Close()

	log.Println("Consumer created successfully, waiting for messages...")

	// Create a channel to handle graceful shutdown
	done := make(chan struct{})

	// Start consuming messages in a goroutine
	go func() {
		consumer.Consume(func(msg rabbitmq.MQMessage) {
			log.Printf("Received message: %v", msg.Data)

			// Convert map to JSON bytes
			jsonData, err := json.Marshal(msg.Data)
			if err != nil {
				log.Printf("Failed to marshal message data: %v", err)
				return
			}

			var mintMsg MintMessage
			if err := json.Unmarshal(jsonData, &mintMsg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return
			}

			log.Printf("Processing mint for wallet: %s", mintMsg.WalletAddress)
			err = MintNFTForWallet(mintMsg.WalletAddress)
			if err != nil {
				log.Printf("Failed to mint NFT: %v", err)
				return
			}

			// Only update KYC status if minting was successful
			log.Println("NFT minted successfully, updating KYC status...")
			err = repo.UpdateKYC(context.Background(), sqlc.KycInfo{
				CitizenID: mintMsg.CitizenID,
				IsActive:  pgtype.Bool{Bool: true, Valid: true},
			})
			if err != nil {
				log.Printf("Failed to update KYC status: %v", err)
				return
			}
			log.Println("KYC status updated successfully")
		})
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Signal the consumer to stop
	close(done)
	log.Println("Mint worker shutting down...")
}

func MintNFTForWallet(wallet string) error {
	log.Printf("Starting mint process for wallet: %s", wallet)

	// 1. Connect to Ethereum client
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return fmt.Errorf("RPC_URL environment variable not set")
	}
	log.Printf("Connecting to RPC: %s", rpcURL)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %v", err)
	}
	defer client.Close()

	// 2. Load private key
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	log.Println("Loading private key...")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	// 3. Prepare contract interaction
	chainID := big.NewInt(2021) // Ronin Saigon testnet chain ID
	log.Println("Parsing ABI...")
	parsedABI, err := abi.JSON(strings.NewReader(kycWalletNFTABI))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	to := common.HexToAddress(wallet)
	log.Printf("Preparing mint transaction for address: %s", to.Hex())
	input, err := parsedABI.Pack("mint", to)
	if err != nil {
		return fmt.Errorf("failed to pack mint function: %v", err)
	}

	// 4. Build transaction
	contractAddress := common.HexToAddress(os.Getenv("KYC_ADDRESS"))
	if contractAddress == (common.Address{}) {
		return fmt.Errorf("KYC_ADDRESS environment variable not set or invalid")
	}
	log.Printf("Contract address: %s", contractAddress.Hex())

	// Get the sender's address
	senderAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	log.Printf("Sender address: %s", senderAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	log.Printf("Current nonce: %d", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}
	log.Printf("Gas price: %s", gasPrice.String())

	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),   // value
		uint64(200_000), // gas limit
		gasPrice,        // gas price
		input,
	)

	log.Println("Signing transaction...")
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 5. Send transaction
	log.Println("Sending transaction...")
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())

	// 6. Wait for transaction confirmation
	log.Println("Waiting for transaction confirmation...")
	for {
		receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
		if err != nil {
			// If receipt is not available yet, continue waiting
			log.Println("Transaction pending...")
			time.Sleep(3 * time.Second)
			continue
		}

		// Transaction has been mined
		if receipt.Status == 1 {
			log.Println("Transaction SUCCESS: NFT minted")
			return nil
		} else {
			log.Println("Transaction FAILED: Mint NFT failed")
			return fmt.Errorf("transaction failed (revert). Check contract logic or parameters")
		}
	}
}
