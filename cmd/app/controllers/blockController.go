package controllers

import (
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