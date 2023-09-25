package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CT77777/Block-Indexer/initializers"
	"github.com/CT77777/Block-Indexer/models"
	"github.com/gin-gonic/gin"
)

// get a limited count of blocks
func GetBlocks(c *gin.Context) {

	var blockHeaders []models.BlockHeader

	count := c.Query("limit")

	countInt, err := strconv.Atoi(count)

	if err != nil {
		c.JSON(400, gin.H{"Error": "Invalid blocks limit"})
		return
	}

	result := initializers.DB.Table("blocks").Order("Number DESC").Limit(countInt).Find(&blockHeaders)

	if result.Error != nil {
		log.Printf("DB issue: %v", result.Error)
		
		c.JSON(500, gin.H{"Error": "Internal server error"})
		return 
	}

	blocks := struct{Blocks []models.BlockHeader `json:"blocks"`}{Blocks: blockHeaders}

	c.JSON(200, blocks)
}

// get the specified block, including all transactions hash
func GetBlockAndTxs(c *gin.Context) {
	var blockAndTx []models.BlockAndTx

	number := c.Param("id")

	_, err := strconv.ParseInt(number, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{"Error":"Invalid block number"})
		return
	}

	// use Joins, interact with DB one time
	result := initializers.DB.Table("blocks").Select("blocks.*", "transactions.hash as Transaction").Where("number = ?", number).Joins("INNER JOIN transactions ON blocks.number = transactions.block_number").Scan(&blockAndTx)

	if result.Error != nil {
		fmt.Printf("Error: %v", result.Error)

		c.JSON(500, gin.H{"Error": "Internal server error"})

		return
	}

	var blockAndTxs models.BlockAndTxs
	
	for index, value := range blockAndTx {
		if index == 0 {
			blockAndTxs.Number = value.Number
			blockAndTxs.Hash = value.Hash
			blockAndTxs.Time = value.Time
			blockAndTxs.Parent_Hash = value.Parent_Hash
		}

		blockAndTxs.Transactions = append(blockAndTxs.Transactions, value.Transaction)
	}

	c.JSON(200, blockAndTxs)
}

// get the specified transaction, include all event logs
func GetTxAndLogs(c *gin.Context) {

	txHash := c.Param("txHash")

	var txAndLog []models.TxAndLog

	result := initializers.DB.Table("transactions").Select("transactions.*", "logs.data").Where("hash = ?", txHash).Joins("INNER JOIN logs ON transactions.id = logs.transaction_id").Scan(&txAndLog)

	if result.Error != nil {
		log.Printf("Error: %v", result.Error)

		c.JSON(500, gin.H{"Error" : "Internal server error"})
		return 
	}

	var txAndLogs models.TxAndLogs

	for index, value := range txAndLog {
		if index == 0 {
			txAndLogs.Hash = value.Hash
			txAndLogs.From = value.From
			txAndLogs.To = value.To
			txAndLogs.Nonce = value.Nonce
			txAndLogs.Value = value.Value
			txAndLogs.Input_Data = value.Input_Data
		}

		var log models.EventLog = models.EventLog{Index : uint8(index), Data : value.Data}

		txAndLogs.Logs = append(txAndLogs.Logs, log)
	}

	c.JSON(200, txAndLogs)
}