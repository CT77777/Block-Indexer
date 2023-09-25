package controllers

import (
	"log"

	"github.com/CT77777/Block-Indexer/initializers"
	"github.com/CT77777/Block-Indexer/models"
	"github.com/gin-gonic/gin"
)

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