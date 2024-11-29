package main

import (
	"WesChess/backend/internal/db"
	"WesChess/backend/internal/handlers"
	"WesChess/backend/internal/matchmaking"
	"WesChess/backend/internal/ws"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := gin.Default()

	router.Static("/static", "../frontend/public")

	router.LoadHTMLGlob("../frontend/public/*.html")

	router.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		ws.HandleConnection(c.Writer, c.Request, roomID)
	})

	router.GET("/", func(c *gin.Context) {
		log.Printf("Redirecting to /home")
		c.Redirect(http.StatusFound, "/home")
	})

	router.GET("/ws-test", func(c *gin.Context) {
		c.File("../frontend/public/ws_test.html")
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	router.GET("/home", func(c *gin.Context) {
		userID, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(302, "/login")
			return
		}

		var username string
		err = db.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			c.HTML(500, "login.html", gin.H{
				"error": "Failed to load user information",
			})
			return
		}

		c.HTML(200, "home.html", gin.H{
			"username": username,
		})
	})

	router.GET("/api/logged-in-users", func(c * gin.Context) {
		users := make([]string, 0, len(handlers.ActiveUsers))
		for _, username := range handlers.ActiveUsers {
			users = append(users, username)
		}
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	})

	router.GET("/api/user", func(c *gin.Context) {
		userIDstr, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var username string
		err = db.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user information"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	router.GET("/api/check-match", func(c *gin.Context) {
		userIDstr, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		matchmaking.MatchedMutex.Lock()
		roomID, matched := matchmaking.MatchedPlayers[userID]
		log.Printf("matchmaking.MatchedPlayers: %v", matchmaking.MatchedPlayers)
		matchmaking.MatchedMutex.Unlock()

		if matched {
			log.Printf("User %d matched with room %d", userID, roomID)
			c.JSON(http.StatusOK, gin.H{"matched": true, "roomID": roomID})
		} else {
			log.Printf("User %d not matched", userID)
			c.JSON(http.StatusOK, gin.H{"matched": false})
		}

	})

	router.GET("/index", func(c *gin.Context) {
		userID, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(302, "/register")
			return
		}

		var username string
		err = db.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			c.HTML(500, "index.html", gin.H{
				"error": "Failed to load user information",
			})
			return
		}

		c.HTML(200, "index.html", gin.H{
			"username": username,
		})
	})

	router.GET("/game/:roomID", func(c *gin.Context) {
		roomIDstr := c.Param("roomID")

		// Optionally validate that the room exists
		// if _, exists := matchmaking.GetActiveGame(roomID); !exists {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		// 	return
		// }
		// get username from user
		roomID, err := strconv.Atoi(roomIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		// print out active games map
		log.Printf("Active games: %v", ws.ActiveGames)
		game, exists := ws.ActiveGames[roomID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}
		userIDstr, err := c.Cookie("user_id")
		if err != nil {
			c.Redirect(302, "/register")
			return
		}
		// use userID to get username
		var username string
		err = db.DB.QueryRow("SELECT username FROM users WHERE id = ?", userIDstr).Scan(&username)
		if err != nil {
			c.HTML(500, "game.html", gin.H{
				"error": "Failed to load user information",
			})
			return
		}
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		
		var playerColor string
		var opponentID int
		if game.WhiteID == userID {
			playerColor = "w"
			opponentID = game.BlackID
		} else if game.BlackID == userID {
			playerColor = "b"
			opponentID = game.WhiteID
		} else {
			playerColor = "spectator"
		}






		// Serve the game page
		c.HTML(http.StatusOK, "game.html", gin.H{
			"roomID":      roomID,
			"username":    username,
			"user_id":     userID,
			"opponentID":    opponentID,
			"whiteID":     game.WhiteID,
			"blackID":     game.BlackID,
			"playerColor": playerColor,
		})
	})

	router.POST("/api/register", handlers.RegisterHandler(db.DB))

	router.POST("/api/play", func(c *gin.Context) {
		userIDstr, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		matchmaking.EnqueuePlayer(int(userID))
		log.Printf("Player %d enqueued", userID)
		p1, p2, room, match_found := matchmaking.MatchPlayers()
		if match_found {
			log.Printf("Match found between players %d and %d in room %d", p1, p2, room)
			c.JSON(http.StatusOK, gin.H{"message": "Match found", "roomID": room})
		} else {
			log.Printf("No match found")
		}
	})

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.POST("/register", handlers.RegisterHandler(db.DB))
	router.POST("/login", handlers.LoginHandler(db.DB))

	port := os.Getenv("PORT") // Get the port from the environment
	if port == "" {
		port = "8080" // Default to 8080 for local development
	}

	log.Println("Server started on :" + port)
	router.Run(":" + port)
}
