package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CT77777/Block-Indexer/db/models"
	"github.com/CT77777/Block-Indexer/initializers"
	"github.com/gin-gonic/gin"
)

// get a limited count of blocks
func GetBlocks(c *gin.Context) {

	var blocks []models.Block

	count := c.Query("limit")

	countInt, err := strconv.Atoi(count)

	if err != nil {
		c.JSON(400, gin.H{"Error": "Invalid blocks limit"})
		return
	}

	result := initializers.DB.Order("Number DESC").Limit(countInt).Find(&blocks)

	if result.Error != nil {
		log.Printf("DB issue: %v", result.Error)
		
		c.JSON(500, gin.H{"Error": "Internal server error"})
		return 
	}

	c.JSON(200, &blocks)
}

type BlockAndTx struct {
	Number uint64
	Hash string
	Time uint64
	Parent_Hash string
	Transaction string
}

type BlockAndTxs struct {
	Number    uint64   `json:"block_num"`
	Hash   string   `json:"block_hash"`
	Time   uint64   `json:"block_time"`
	Parent_Hash  string   `json:"parent_hash"`
	Transactions []string `json:"transactions"`
}

// get the specified block, including all transactions hash
func GetBlockAndTxs(c *gin.Context) {
	var blockAndTx []BlockAndTx

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

	var blockAndTxs BlockAndTxs
	
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