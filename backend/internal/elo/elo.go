package elo

import (
	"database/sql"
	"log"
)

func IncrementWins(userID int, db *sql.DB) {
	_, err := db.Exec("UPDATE users SET wins = wins + 1 WHERE id = ?", userID)
	if err != nil {
		log.Printf("Failed to increment wins for user %d: %v", userID, err)
	}
}

func IncrementLosses(userID int, db *sql.DB) {
	_, err := db.Exec("UPDATE users SET losses = losses + 1 WHERE id = ?", userID)
	if err != nil {
		log.Printf("Failed to increment losses for user %d: %v", userID, err)
	}
}

func IncrementDraws(userID int, db *sql.DB) {
	_, err := db.Exec("UPDATE users SET draws = draws + 1 WHERE id = ?", userID)
	if err != nil {
		log.Printf("Failed to increment draws for user %d: %v", userID, err)
	}
}