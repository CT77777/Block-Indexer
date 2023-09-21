package main

import (
	"github.com/CT77777/Block-Indexer/poc/initializers"
	"github.com/CT77777/Block-Indexer/poc/models"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
}