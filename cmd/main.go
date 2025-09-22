package main

import (
	"github.com/gin-gonic/gin"
)

type SignRequest struct {
	Tg_name  string `json:"tg_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is running"})
	})

	r.POST("/sing-up", handleSingUp)
	r.POST("/sing-in", handlerSingIn)
	r.Run()
}

func handleSingUp(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "successfully signed up"})
}

func handlerSingIn(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"jwt": "hahahahhahaha"})
}
