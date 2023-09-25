package initializers

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *ethclient.Client

func InitClient() {
	var err error
	rpcEndpoint := os.Getenv("RPC_ENDPOINT")

	EthClient, err = ethclient.Dial(rpcEndpoint)

	if err != nil {
		log.Fatal("Failed to establish connection to Ethereum node", err)
	}
}