package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	rpcEndpoint := os.Getenv("RPC_ENDPOINT")

	client, err := ethclient.Dial(rpcEndpoint)

	if err != nil {
		log.Fatal(err)
	}
	
	// get detailed data of the block
	block , err := client.BlockByNumber(context.TODO(),nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Block Number: %v\n", block.Number())
	log.Printf("Block Hash: %v\n", block.Hash().Hex())
	log.Printf("Block Time: %v\n", block.Time())
	log.Printf("Parent Hash: %v\n", block.ParentHash().Hex())
	log.Printf("Txs: %v\n", block.Transactions())
	log.Printf("Gas Limit: %v\n", block.GasLimit())
	log.Printf("Gas Used: %v\n", block.GasUsed())
	log.Printf("Number of Transactions: %d\n", len(block.Transactions()))
    log.Println("Transaction Hashes:")

    for _, tx := range block.Transactions() {
        log.Printf("Tx Hash: %v\n", tx.Hash().Hex())
    }

	//  // Subscribe to new block headers
	//  headers := make(chan *types.Header)
	//  sub, err := client.SubscribeNewHead(context.TODO(), headers) // only support websocket RPC endpoint
	//  if err != nil {
	// 	 log.Fatalf("Failed to subscribe to new block headers: %v", err)
	//  }
	//  defer sub.Unsubscribe()
 
	//  for {
	// 	 select {
	// 	 case <-sub.Err():
	// 		 log.Println("Block header subscription closed.")
	// 		 return
	// 	 case header := <-headers:
	// 		 log.Printf("New block received: BlockNumber %v, Hash %s\n", header.Number, header.Hash().Hex())
	// 		 // Handle the new block header as needed
	// 	 }
	//  }


	// using http request to polling the latest block
	 // Initialize the last known block number
	 currentBlockNumber, err := client.BlockNumber(context.Background())
		 if err != nil {
			 log.Printf("Failed to retrieve the latest block number: %v", err)
		 }
	 lastBlockNumber := new(big.Int).SetUint64(currentBlockNumber)

	 // Poll for new blocks at regular intervals
	 pollInterval := 5 * time.Second // Adjust the interval as needed
	 for {
		 // Fetch the latest block number
		 currentBlockNumber, err := client.BlockNumber(context.Background())
		 if err != nil {
			 log.Printf("Failed to retrieve the latest block number: %v", err)
			 continue
		 }
 
		 // Check for new blocks starting from the last known block number
		 currentBlockBigInt := new(big.Int).SetUint64(currentBlockNumber) // Convert to *big.Int
		 for blockNum := new(big.Int).Set(lastBlockNumber).Add(big.NewInt(1), lastBlockNumber); blockNum.Cmp(currentBlockBigInt) < 0; blockNum.Add(blockNum, big.NewInt(1)) {
			 block, err := client.BlockByNumber(context.Background(), blockNum)
			 if err != nil {
				 log.Printf("Failed to retrieve block %s: %v", blockNum.String(), err)
				 continue
			 }
 
			 // Handle the new block
			 log.Printf("New block received: BlockNumber %v, Hash %s\n", block.Number(), block.Hash().Hex())
 
			 // Update the last known block number
			 lastBlockNumber.Set(blockNum)
		 }
 
		 time.Sleep(pollInterval)
	 }

}

// func (*Client) HeaderByNumber -> get the block data, not include transactions hash
// func (*Client) BlockByNumber  -> get the block data, include transactions hash
// func (*Client) TransactionByHash -> get the transaction data, not include event logs
// func (*Client) TransactionReceipt -> get the transaction data, include event logs