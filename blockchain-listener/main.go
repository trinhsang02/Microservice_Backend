package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"blockchain-listener/rabbitmq"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractAddress = "0xf7247772787Dcf46FDfA98CAC15382fF98eaE225"
const RPC_URL = "wss://ronin-saigon.g.alchemy.com/v2/sA-DN7hd8Jk7m34FB5GH2Zt_wb0u_Tud"

// ABI for the Deposit and Withdrawal events
const contractABI = `[
  {"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"commitment","type":"bytes32"},{"indexed":true,"internalType":"address","name":"depositor","type":"address"},{"indexed":false,"internalType":"uint32","name":"leafIndex","type":"uint32"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"Deposit","type":"event"},
  {"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"bytes32","name":"nullifierHash","type":"bytes32"},{"indexed":true,"internalType":"address","name":"relayer","type":"address"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"Withdrawal","type":"event"}
]`

func main() {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	queueName := "blockchain_events"

	producer, err := rabbitmq.NewProducer(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum RPC: %v", err)
	}
	log.Println("Connected to Ethereum RPC")

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contract := common.HexToAddress(contractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
	}

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
