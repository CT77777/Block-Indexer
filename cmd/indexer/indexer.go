package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/CT77777/Block-Indexer/initializers"
	"github.com/CT77777/Block-Indexer/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// initialize prerequisite
func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.InitClient()
}

// worker for consuming jobs in channel
func worker(id int, jobs <-chan int, results chan<- models.Block) {
	for job := range jobs {
		blockNumber := big.NewInt(int64(job))

		fmt.Printf("Block# %v worker_%v\n", blockNumber, id)
		
		block, err := initializers.EthClient.BlockByNumber(context.TODO(), blockNumber)

		if err != nil {
			log.Printf(`Failed to fetch Block #%v`, job)
			continue;
		}

		// get transaction data by transaction hash
		var transactions []models.Transaction

		for _, txHash := range block.Transactions() {
			hash := txHash.Hash().Hex()

			tx , isPending ,err := initializers.EthClient.TransactionByHash(context.TODO(), common.HexToHash(hash))

			if err != nil {
				log.Printf(`Failed to fetch Transaction %v`, hash)
				continue
			}

			receipt, err := initializers.EthClient.TransactionReceipt(context.TODO(), common.HexToHash(hash))

			if err != nil {
				log.Printf(`Failed to fetch event logs of Transaction %v`, hash)
				continue
			}

			if isPending != true {
				var logs []models.Log

				for _, log := range receipt.Logs {
					logData := models.Log{Address: log.Address.Hex(), Data: `0x`+hex.EncodeToString(log.Data)}

					for index, topic := range log.Topics {
						topicHex := topic.Hex()

						switch index {
							case 0:
								logData.Topic_0 = topicHex
							case 1:
								logData.Topic_1 = topicHex
							case 2:
								logData.Topic_2 = topicHex
							case 3:
								logData.Topic_3 = topicHex
						}
					}

					logs = append(logs, logData)
				}

				from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)

				if err != nil {
					log.Println(`Failed to get sender address`, err)
					continue
				}

				var transactionData models.Transaction

				if tx.To() != nil {
					transactionData = models.Transaction{Hash: tx.Hash().Hex(), From: from.Hex(), To: tx.To().Hex(), Nonce: tx.Nonce(), Value: tx.Value().String(), Input_Data: `0x`+hex.EncodeToString(tx.Data()), Logs: logs}
				} else {
					transactionData = models.Transaction{Hash: tx.Hash().Hex(), From: from.Hex(), Nonce: tx.Nonce(), Value: tx.Value().String(), Input_Data: `0x`+hex.EncodeToString(tx.Data()), Logs: logs}
				}

				transactions = append(transactions, transactionData)
			} else {
				log.Printf(`Transaction %v is pending`, hash)
				continue
			}
		}

		blockData := models.Block{Number: block.Number().Uint64(), Hash: block.Hash().Hex(), Time: block.Time(), Parent_Hash: block.ParentHash().Hex(), Transactions: transactions}

		results <- blockData
	}
}

// create workers by goroutine, and listen results channel to store blocks data into DB
func scanBlocks(start int, end int, workerCount int, batchInsertCount int) {
	// create jobs channel & results channel
	jobs := make(chan int, end-start+1)
	results := make(chan models.Block, end-start+1)

	// open workers
	for i:=1; i <= workerCount; i++ {
		go worker(i, jobs, results)
	}

	// push jobs into jobs channel
	for i:=start; i <= end; i++ {
		jobs <- i
	}
	close(jobs)

	// listen block data in results channel, and store them into database
	var blockDatas []models.Block 
	
	for i:=1; i <= end-start+1; i++ {
		blockData := <- results
	
		blockDatas = append(blockDatas, blockData)
		
		// Batch insert for reducing I/O times
		if len(blockDatas) == batchInsertCount {
			result := initializers.DB.Create(&blockDatas)

			if result.Error != nil {
				log.Printf("Failed to insert block #%v data: %v", blockData.Number, result.Error)
			}

			blockDatas = blockDatas[:0]
		} 
	}

	log.Printf("Scanning blocks done. From block #%v to #%v\n", start, end)
 }

// main execution function
func main() {
	// parameters for scanning historical blocks
	start := 0 // start block number
	end := 33585920 // end block number
	workerCountOld := 5 // worker count to scan blocks in parallel for historical blocks
	batchInsertCountOld := 10 // total block data count will be insert into DB every I/O for historical blocks

	// parameters for keeping scanning the latest blocks
	workerCountNew := 1 // worker count to scan blocks in parallel for new blocks
	batchInsertCountNew := 1 // total block data count will be insert into DB every I/O for new blocks
	fetchInterval := 5 * time.Second // periodically fetching interval

	// scan historical blocks
	go scanBlocks(start, end, workerCountOld, batchInsertCountOld)

	// keep scanning new blocks
	lastBlockNum := end

	for {
		fmt.Println("Scanning new blocks...")
		currBlockNum, err := initializers.EthClient.BlockNumber(context.TODO())

		if err != nil {
			log.Fatal("Failed to fetch the latest block number")
		}

		scanBlocks(int(lastBlockNum+1), int(currBlockNum), workerCountNew, batchInsertCountNew)

		lastBlockNum = int(currBlockNum)

		time.Sleep(fetchInterval)
	}
}