package handlers

import (
	"fmt"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"net/http"
    "strconv"
    "log"
)

var ActiveUsers = make(map[string]string)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
            log.Println("Error hashing password")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password, wins, draws, losses, elo) VALUES (?, ?, 0, 0, 0, 1200)", req.Username, hashedPassword)
		if err != nil {
            log.Println("Error inserting user into db")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user"})
			return
		}

		c.Redirect(http.StatusFound, "/login")
		// c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        // Bind JSON payload to the struct
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
            return
        }

        var userID int
        var storedHash string

        // Fetch the user ID and hashed password from the database
        err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &storedHash)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        // Verify the password
        if !CheckPasswordHash(req.Password, storedHash) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        // Set the cookie with the user ID
        c.SetCookie(
            "user_id",              // Name
            fmt.Sprintf("%d", userID), // Value (convert int to string)
            3600,                   // MaxAge in seconds (1 hour in this case)
            "/",                    // Path
            "",                     // Domain (empty means the current domain)
            false,                  // Secure (true if using HTTPS)
            true,                   // HttpOnly (prevents access via JavaScript)
        )

        // add to active users global var
        ActiveUsers[strconv.Itoa(userID)] = req.Username

        c.Redirect(http.StatusFound, "/home")
    }
}



// func LoginHandler(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var req struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 		}

// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 			return
// 		}

// 		var storedHash string
// 		err := db.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedHash)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 			return
// 		}

// 		if !CheckPasswordHash(req.Password, storedHash) {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
// 	}
// }
