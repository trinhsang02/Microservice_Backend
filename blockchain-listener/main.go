package main

import (
	"context"
	"encoding/hex"
	// "encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/yourusername/yourrepo/db/sqlc"
	// "github.com/yourusername/yourrepo/mq/rabbitmq"
)

// ABI for the Deposit and Withdrawal events
const contractABI = `[
  {"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"commitment","type":"bytes32"},{"indexed":true,"internalType":"address","name":"depositor","type":"address"},{"indexed":false,"internalType":"uint32","name":"leafIndex","type":"uint32"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"Deposit","type":"event"},
  {"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"bytes32","name":"nullifierHash","type":"bytes32"},{"indexed":true,"internalType":"address","name":"relayer","type":"address"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"Withdrawal","type":"event"}
]`

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	mixer_0_1 := os.Getenv("MIXER_0_1")
	mixer_1 := os.Getenv("MIXER_1")
	mixer_10 := os.Getenv("MIXER_10")
	mixer_100 := os.Getenv("MIXER_100")
	rpcURL := os.Getenv("RONIN_WEBSOCKET_URL")
	startBlock := os.Getenv("DEPLOYMENT_BLOCK")

	log.Printf("Mixer 0.1 Address: %s", mixer_0_1)
	log.Printf("Mixer 1 Address: %s", mixer_1)
	log.Printf("Mixer 10 Address: %s", mixer_10)
	log.Printf("Mixer 100 Address: %s", mixer_100)
	log.Printf("RPC URL: %s", rpcURL)

	// Test RabbitMQ connection
	// log.Printf("Testing RabbitMQ connection...")
	// testProducer, err := rabbitmq.NewProducer(rabbitmqURL, "blockchain_exchange", "topic")
	// if err != nil {
	// 	log.Printf("Failed to connect to RabbitMQ: %v", err)
	// } else {
	// 	log.Printf("Successfully connected to RabbitMQ")
	// 	defer testProducer.Close()

	// 	// Test publish
	// 	log.Printf("Testing message publish...")
	// 	err = testProducer.PublishStruct("listener.deposit", sqlc.Deposit{
	// 		Commitment:  pgtype.Text{String: "0x123", Valid: true},
	// 		Depositor:   pgtype.Text{String: "0x456", Valid: true},
	// 		LeafIndex:   pgtype.Int4{Int32: 0, Valid: true},
	// 		Timestamp:   pgtype.Numeric{Int: big.NewInt(1643723400), Valid: true},
	// 		TxHash:      pgtype.Text{String: "0x789", Valid: true},
	// 		BlockNumber: pgtype.Int4{Int32: 123456, Valid: true},
	// 		ChainID:     pgtype.Int4{Int32: 2021, Valid: true},
	// 	})
	// 	if err != nil {
	// 		log.Printf("Failed to publish test message: %v", err)
	// 	} else {
	// 		log.Printf("Successfully published test message")
	// 	}
	// }

	if mixer_0_1 == "" || mixer_1 == "" || mixer_10 == "" || mixer_100 == "" || rpcURL == "" || startBlock == "" {
		log.Fatalf("Missing environment variables")
	}

	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	queries := sqlc.New(conn)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum RPC: %v", err)
	}
	log.Println("Connected to Ethereum RPC")

	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Convert all contract addresses
	mixer0_1Contract := common.HexToAddress(mixer_0_1)
	mixer1Contract := common.HexToAddress(mixer_1)
	mixer10Contract := common.HexToAddress(mixer_10)
	mixer100Contract := common.HexToAddress(mixer_100)

	// Start sync up in a goroutine
	go func() {
		log.Println("Starting missed events sync in background...")
		parsedABI, err := abi.JSON(strings.NewReader(contractABI))
		if err != nil {
			log.Fatalf("Failed to parse contract ABI: %v", err)
		}
		contractAddresses := []common.Address{mixer0_1Contract, mixer1Contract, mixer10Contract, mixer100Contract}
		startBlock, err := strconv.ParseUint(startBlock, 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse start block: %v", err)
		}
		blockChunkSize := uint64(499)
		SyncUpEvents(context.Background(), parsedABI, contractAddresses, queries, startBlock, blockChunkSize)
	}()

	query := ethereum.FilterQuery{
		Addresses: []common.Address{mixer0_1Contract, mixer1Contract, mixer10Contract, mixer100Contract},
	}

	log.Printf("FilterQuery: %+v", query)

	logs := make(chan types.Log)
	ctx := context.Background()
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to contract events: %v", err)
		log.Println("Listening for Deposit and Withdrawal events...")
		parsedABI, _ := abi.JSON(strings.NewReader(contractABI))

		for {
			select {
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case vLog := <-logs:

				event, err := parsedABI.EventByID(vLog.Topics[0])
				if err != nil {
					log.Printf("Unknown event: %v", err)
					continue
				}

				if event.Name == "Deposit" {
					data := make(map[string]interface{})
					err = parsedABI.UnpackIntoMap(data, event.Name, vLog.Data)

					if err != nil {
						log.Printf("Failed to unpack event: %v", err)
						continue
					}
					// Extract indexed parameters
					commitment := common.BytesToHash(vLog.Topics[1][:]).Hex()
					depositor := common.BytesToAddress(vLog.Topics[2][:]).Hex()

					// Extract non-indexed parameters
					leafIndex := uint32(0)
					if leafIndexVal, ok := data["leafIndex"].(uint32); ok {
						leafIndex = leafIndexVal
					}

					// Extract the timestamp as *big.Int from the event data
					if timestampVal, ok := data["timestamp"].(*big.Int); ok {
						deposit, err := queries.CreateDeposit(context.Background(), sqlc.CreateDepositParams{
							ContractAddress: pgtype.Text{String: vLog.Address.Hex(), Valid: true},
							Commitment:      pgtype.Text{String: commitment, Valid: true},
							Depositor:       pgtype.Text{String: depositor, Valid: true},
							LeafIndex:       pgtype.Int4{Int32: int32(leafIndex), Valid: true},
							Timestamp:       pgtype.Numeric{Int: timestampVal, Valid: true},
							TxHash:          pgtype.Text{String: vLog.TxHash.Hex(), Valid: true},
							BlockNumber:     pgtype.Int4{Int32: int32(vLog.BlockNumber), Valid: true},
							ChainID:         pgtype.Int4{Int32: int32(2021), Valid: true},
						})
						if err != nil {
							log.Printf("Failed to insert deposit: %v", err)
						} else {
							if deposit.ID != 0 {
								log.Printf("Deposit event stored: commitment=%s, depositor=%s, timestamp=%s, txHash=%s", commitment, depositor, timestampVal.String(), vLog.TxHash.Hex())
							}
						}
					}
				} else {
					// Extract indexed parameter
					relayer := common.BytesToAddress(vLog.Topics[1][:]).Hex()

					// For Withdrawal events, data contains:
					// [0:32]   - recipient (address)
					// [32:64]  - nullifierHash (bytes32)
					// [64:96]  - fee (uint256)
					if len(vLog.Data) < 96 {
						log.Printf("Invalid data length for Withdrawal event: %d", len(vLog.Data))
						continue
					}

					// Extract recipient (first 32 bytes, but only last 20 bytes are the address)
					recipient := common.BytesToAddress(vLog.Data[12:32]).Hex()

					// Extract nullifierHash (next 32 bytes)
					nullifier := "0x" + hex.EncodeToString(vLog.Data[32:64])

					// Extract fee (last 32 bytes)
					fee := new(big.Int).SetBytes(vLog.Data[64:96])

					log.Printf("Withdrawal data extracted - recipient: %s, nullifier: %s, relayer: %s, fee: %s",
						recipient, nullifier, relayer, fee.String())

					// Store in database
					withdrawal, err := queries.CreateWithdrawal(context.Background(), sqlc.CreateWithdrawalParams{
						ContractAddress: pgtype.Text{String: vLog.Address.Hex(), Valid: true},
						NullifierHash:   pgtype.Text{String: nullifier, Valid: true},
						Recipient:       pgtype.Text{String: recipient, Valid: true},
						Relayer:         pgtype.Text{String: relayer, Valid: true},
						Fee:             pgtype.Numeric{Int: fee, Valid: true},
						Timestamp:       pgtype.Numeric{Int: big.NewInt(int64(vLog.BlockNumber)), Valid: true},
						TxHash:          pgtype.Text{String: vLog.TxHash.Hex(), Valid: true},
						BlockNumber:     pgtype.Int4{Int32: int32(vLog.BlockNumber), Valid: true},
						ChainID:         pgtype.Int4{Int32: int32(2021), Valid: true},
					})
					if err != nil {
						log.Printf("Failed to insert withdrawal: %v", err)
					} else {
						if withdrawal.ID != 0 {
							log.Printf("Withdrawal event stored: nullifierHash=%s, recipient=%s, relayer=%s, fee=%s, txHash=%s",
								nullifier, recipient, relayer, fee.String(), vLog.TxHash.Hex())
						}
					}
				}
			}
			// TODO: push to rabbitmq

			// log.Printf("Setting up RabbitMQ consumer...")
			// consumer, err := rabbitmq.NewConsumer(rabbitmqURL, "blockchain_exchange", "topic", "blockchain_listener_queue", []string{"listener.*"})
			// if err != nil {
			// 	log.Printf("Failed to initialize RabbitMQ consumer: %v", err)
			// } else {
			// 	log.Printf("Successfully set up RabbitMQ consumer")
			// 	defer consumer.Close()

			// 	err = consumer.Consume(func(msg rabbitmq.MQMessage) {
			// 		log.Printf("Received message of type: %s", msg.Type)
			// 		switch msg.Type {
			// 		case "listener.deposit":
			// 			// Xử lý deposit
			// 			var deposit sqlc.Deposit
			// 			b, _ := json.Marshal(msg.Data)
			// 			_ = json.Unmarshal(b, &deposit)
			// 			log.Printf("Received deposit: %+v", deposit)
			// 		case "listener.withdraw":
			// 			// Xử lý withdraw
			// 			var withdraw sqlc.Withdrawal
			// 			b, _ := json.Marshal(msg.Data)
			// 			_ = json.Unmarshal(b, &withdraw)
			// 			log.Printf("Received withdraw: %+v", withdraw)
			// 		default:
			// 			log.Printf("Unknown message type: %s", msg.Type)
			// 		}
			// 	})
			// 	if err != nil {
			// 		log.Printf("Failed to consume messages: %v", err)
			// 	} else {
			// 		log.Printf("Successfully started consuming messages")
			// 	}
			// }
		}
	}
}
