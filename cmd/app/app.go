package main

import (
	"github.com/CT77777/Block-Indexer/cmd/app/controllers"
	"github.com/CT77777/Block-Indexer/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	app := gin.Default()

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Status": true})
	})

	app.GET("/blocks", controllers.GetBlocks)

	app.Run()
}