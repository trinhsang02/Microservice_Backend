package main

import (
	"blockchain-listener/rabbitmq"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
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

	log.Printf("RabbitMQ URL: %s", rabbitmqURL)
	log.Printf("Mixer 0.1 Address: %s", mixer_0_1)
	log.Printf("Mixer 1 Address: %s", mixer_1)
	log.Printf("Mixer 10 Address: %s", mixer_10)
	log.Printf("Mixer 100 Address: %s", mixer_100)
	log.Printf("RPC URL: %s", rpcURL)

	if mixer_0_1 == "" || mixer_1 == "" || mixer_10 == "" || mixer_100 == "" {
		log.Fatalf("Missing environment variables")
	}

	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	queueName := "blockchain_events"

	producer, err := rabbitmq.NewProducer(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

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

	client, err := ethclient.Dial("https://ronin-saigon.g.alchemy.com/v2/o18cYN4bRHQDLeD10ewVI17551htkNKg")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum RPC: %v", err)
	}
	log.Println("Connected to Ethereum RPC")

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := common.HexToAddress(mixer_0_1)

	// Get the latest block number
	latestBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	// Start from a specific block or use a reasonable starting point
	startBlock := uint64(37659740)
	blockChunkSize := uint64(499)

	for currentBlock := startBlock; currentBlock <= latestBlock; currentBlock += blockChunkSize {
		endBlock := currentBlock + blockChunkSize
		if endBlock > latestBlock {
			endBlock = latestBlock
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(currentBlock)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: []common.Address{contract},
			Topics: [][]common.Hash{{
				parsedABI.Events["Deposit"].ID,
			}},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("Error filtering logs: %v", err)
			continue
		}

		for _, vLog := range logs {
			// Process Deposit event
			data := make(map[string]interface{})
			err = parsedABI.UnpackIntoMap(data, "Deposit", vLog.Data)
			if err != nil {
				log.Printf("Failed to unpack Deposit event: %v", err)
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
				// Store in database
				_, err = queries.CreateDeposit(context.Background(), sqlc.CreateDepositParams{
					ContractAddress: pgtype.Text{String: contract.Hex(), Valid: true},
					Commitment:      pgtype.Text{String: commitment, Valid: true},
					Depositor:       pgtype.Text{String: depositor, Valid: true},
					LeafIndex:       pgtype.Int4{Int32: int32(leafIndex), Valid: true},
					Timestamp:       pgtype.Numeric{Int: timestampVal, Valid: true},
					TxHash:          pgtype.Text{String: vLog.TxHash.Hex(), Valid: true},
				})
				if err != nil {
					log.Printf("Failed to insert deposit: %v", err)
				} else {
					log.Printf("Deposit event stored: commitment=%s, depositor=%s, timestamp=%s, txHash=%s", commitment, depositor, timestampVal.String(), vLog.TxHash.Hex())
				}
			}
		}

		log.Printf("Processed blocks %d to %d, found %d deposit events", currentBlock, endBlock, len(logs))
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
		FromBlock: big.NewInt(37511375),
	}

	log.Printf("FilterQuery: %+v", query)

	logs := make(chan types.Log)
	ctx := context.Background()
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to contract events: %v", err)
	}
	log.Println("Listening for Deposit and Withdrawal events...")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error: %v", err)
			return
		case vLog := <-logs:
			// Identify event type
			event, err := parsedABI.EventByID(vLog.Topics[0])
			if err != nil {
				log.Printf("Unknown event: %v", err)
				continue
			}
			data := make(map[string]interface{})
			err = parsedABI.UnpackIntoMap(data, event.Name, vLog.Data)
			if err != nil {
				log.Printf("Failed to unpack event: %v", err)
				continue
			}
			// Add indexed fields
			for i, input := range event.Inputs {
				if input.Indexed {
					if len(vLog.Topics) > i+1 {
						if input.Type.String() == "address" {
							data[input.Name] = common.HexToAddress(vLog.Topics[i+1].Hex()).Hex()
						} else {
							data[input.Name] = vLog.Topics[i+1].Hex()
						}
					}
				}
			}
			// Prepare message
			msg := event.Name + ": "
			for k, v := range data {
				msg += k + "=" + toString(v) + ", "
			}
			log.Printf("Event: %s", msg)
			err = producer.Publish(msg)
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
			}
		}
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case *big.Int:
		return val.String()
	case common.Address:
		return val.Hex()
	case [32]byte:
		return "0x" + hex.EncodeToString(val[:])
	case uint32:
		return fmt.Sprintf("%d", val)
	case uint64:
		return fmt.Sprintf("%d", val)
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return "<unknown>"
	}
}
