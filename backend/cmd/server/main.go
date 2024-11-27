package main

import (
	"WesChess/backend/internal/db"
	"WesChess/backend/internal/handlers"
	"WesChess/backend/internal/ws"
	"WesChess/backend/internal/matchmaking"
	"log"
	"strconv"
	"net/http"

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
		c.Redirect(302, "/register")
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

		c.HTML(200, "home.html", gin.H{
			"username": username,
		})
	})

	router.GET("/api/user", func (c *gin.Context) {
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
		roomID , err := strconv.Atoi(roomIDstr)
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
		if game.WhiteID == userID {
			playerColor = "white"
		} else if game.BlackID == userID {
			playerColor = "black"
		} else {
			playerColor = "spectator"
		}

		// Serve the game page
		c.HTML(http.StatusOK, "game.html", gin.H{
			"roomID": roomID,
			"username": username,
			"user_id": userID,
			"whiteID": game.WhiteID,
			"blackID": game.BlackID,
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
	log.Println("Server started on :8080")

	router.Run(":8080")
}

