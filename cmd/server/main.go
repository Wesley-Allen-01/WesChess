package main

import (
	"WesChess/internal/db"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	log.Println("Server started on :8080")

	router.Run(":8080")
}

