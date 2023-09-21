package main

import (
	"github.com/CT77777/Block-Indexer/poc/controllers"
	"github.com/CT77777/Block-Indexer/poc/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
}

func main() {
	app := gin.Default()

	app.POST("/post", controllers.PostCreate)
	app.GET("/post", controllers.PostGetAll)
	app.GET("/post/:id", controllers.PostGetOne)
	app.PUT("/post/:id", controllers.PostUpdate)
	app.DELETE("/post/:id", controllers.PostDelete)

	app.Run()
}