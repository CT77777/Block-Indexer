package main

import (
	"github.com/CT77777/Block-Indexer/db/models"
	"github.com/CT77777/Block-Indexer/initializers"
)


func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	block := models.Block{}
	transaction := models.Transaction{}
	log := models.Log{}

	initializers.DB.AutoMigrate(block, transaction, log)
}