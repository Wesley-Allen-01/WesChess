package db

import (
	"database/sql"
	"log"
	"WesChess/backend/internal/handlers"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./data/chess_game.db")
	if err != nil {
		log.Fatalf("i failed and im sorry: %v", err)
	}

	createTables()
}

func createTables() {
	usersTable := `
		DROP TABLE IF EXISTS users;
		CREATE TABLE IF NOT EXISTS users (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		username TEXT NOT NULL,
    		password TEXT NOT NULL,
			wins INTEGER NOT NULL,
			draws INTEGER NOT NULL,
			losses INTEGER NOT NULL,
			elo INTEGER NOT NULL
		);
	`
	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	
	hashedPassword, err := handlers.HashPassword("wes")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	_, err = DB.Exec(`INSERT INTO users (username, password, wins, draws, losses, elo) VALUES ("wes", ?, 0, 0, 0, 2890);`, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to insert user into users table: %v", err)
	}
	log.Println("Successfully created users table and added wes user")
	gamesTable := `
		DROP TABLE IF EXISTS games;
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player_white_id INTEGER NOT NULL,
			player_black_id INTEGER NOT NULL,
			board_state TEXT NOT NULL,
			status TEXT NOT NULL,
			turn TEXT NOT NULL,
			result TEXT
		);
	`

	_, err = DB.Exec(gamesTable)
	if err != nil {
		log.Fatalf("Failed to create games table: %v", err)
	}

	log.Println("Successfully created tables")

}

