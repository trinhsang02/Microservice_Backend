package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourusername/yourrepo/db/sqlc"
)

func SyncUpEvents(
	ctx context.Context,
	parsedABI abi.ABI,
	contractAddresses []common.Address,
	queries *sqlc.Queries,
	startBlock uint64,
	blockChunkSize uint64,
) {
	// Get RPC URL from environment
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatalf("RPC_URL environment variable is not set")
	}

	// Initialize Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum RPC: %v", err)
	}
	log.Printf("Connected to Ethereum RPC at %s for sync up", rpcURL)
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Printf("Failed to get latest block number: %v", err)
		return
	}

	lastSyncedBlock, err := queries.GetLatestDepositSyncedBlock(ctx, sqlc.GetLatestDepositSyncedBlockParams{
		ContractAddress: pgtype.Text{String: contractAddresses[0].Hex(), Valid: true},
		ChainID:         pgtype.Int4{Int32: int32(2021), Valid: true},
	})

	if err == nil {
		if v, ok := lastSyncedBlock.(int64); ok {
			if uint64(v) > startBlock {
				startBlock = uint64(v)
			}
		}
	}

	for currentBlock := startBlock; currentBlock <= latestBlock; currentBlock += blockChunkSize {
		endBlock := currentBlock + blockChunkSize
		if endBlock > latestBlock {
			endBlock = latestBlock
		}

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(currentBlock)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: contractAddresses,
		}

		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			log.Printf("Error filtering logs: %v", err)
			continue
		}

		if len(logs) == 0 {
			continue
		}

		for _, vLog := range logs {
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
				commitment := common.BytesToHash(vLog.Topics[1][:]).Hex()
				depositor := common.BytesToAddress(vLog.Topics[2][:]).Hex()
				leafIndex := uint32(0)
				if leafIndexVal, ok := data["leafIndex"].(uint32); ok {
					leafIndex = leafIndexVal
				}
				if timestampVal, ok := data["timestamp"].(*big.Int); ok {
					_, err = queries.CreateDeposit(ctx, sqlc.CreateDepositParams{
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
						log.Printf("Deposit event stored: commitment=%s, depositor=%s, timestamp=%s, txHash=%s", commitment, depositor, timestampVal.String(), vLog.TxHash.Hex())
					}
				}
			} else if event.Name == "Withdrawal" {
				relayer := common.BytesToAddress(vLog.Topics[1][:]).Hex()
				if len(vLog.Data) < 96 {
					log.Printf("Invalid data length for Withdrawal event: %d", len(vLog.Data))
					continue
				}
				recipient := common.BytesToAddress(vLog.Data[12:32]).Hex()
				nullifier := "0x" + hex.EncodeToString(vLog.Data[32:64])
				fee := new(big.Int).SetBytes(vLog.Data[64:96])
				_, err = queries.CreateWithdrawal(ctx, sqlc.CreateWithdrawalParams{
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
					log.Printf("Withdrawal event stored: nullifierHash=%s, recipient=%s, relayer=%s, fee=%s, txHash=%s", nullifier, recipient, relayer, fee.String(), vLog.TxHash.Hex())
				}
			}
		}
	}

	log.Printf("Sync up completed for blocks %d to %d", startBlock, latestBlock)
}
