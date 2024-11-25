package main

import (
	"WesChess/backend/internal/db"
	"WesChess/backend/internal/handlers"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := gin.Default()


	router.Static("/static", "frontend/public")

	router.GET("/register", func(c *gin.Context) {
		c.File("/Users/wesleyallen/programming/WesChess/frontend/public/register.html")
	})

	router.POST("/api/register", handlers.RegisterHandler(db.DB))

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.POST("/register", handlers.RegisterHandler(db.DB))
	router.POST("/login", handlers.LoginHandler(db.DB))
	log.Println("Server started on :8080")

	router.Run(":8080")
}

