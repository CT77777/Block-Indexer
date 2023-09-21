package controllers

import (
	"github.com/CT77777/Block-Indexer/poc/initializers"
	"github.com/CT77777/Block-Indexer/poc/models"
	"github.com/gin-gonic/gin"
)

func PostCreate(ctx *gin.Context) {
	// Get data of req body
	var body struct  {
		Title string
		Body string
	}

	ctx.Bind(&body)

	// Create a post
	post := models.Post{Title: body.Title, Body: body.Body,}

	result := initializers.DB.Create(&post) // pass pointer of data to Create

	if result.Error != nil {
		ctx.Status(400)
		return
	}

	// Return it
	ctx.JSON(200, gin.H{"Post" : post})
}

func PostGetAll(ctx *gin.Context) {
	// get the post
	var posts []models.Post
	initializers.DB.Find(&posts)

	// respond with it
	ctx.JSON(200, gin.H{"Posts" : posts})
}

func PostGetOne(ctx *gin.Context) {
	// get the id off url

	id := ctx.Param("id")

	// get the specified post
	var post models.Post
	initializers.DB.First(&post, id)

	// respond with it
	ctx.JSON(200, gin.H{"Post" : post})
}

func PostUpdate(ctx *gin.Context) {
	// get the id off the url
	id := ctx.Param("id")

	// get the data off req body
	var body struct {
		Title string
		Body string
	}

	ctx.Bind(&body)

	// find the post which is ready to be updated
	var post models.Post
	initializers.DB.First(&post, id)

	// update it 
	initializers.DB.Model(&post).Updates(models.Post{
		Title : body.Title, Body: body.Body,
	})

	// Respond it 
	ctx.JSON(200, gin.H{"Updated Post" : post})
}

func PostDelete(ctx *gin.Context) {
	// get the id off the url
	id := ctx.Param("id")

	// delete the specified post
	var post models.Post
	initializers.DB.Delete(&post, id)

	// respond it
	ctx.Status(200)
}