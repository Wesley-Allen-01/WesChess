package matchmaking

import (
	"WesChess/backend/internal/ws"
	"log"
	"sync"
)

const InitialBoardState = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Create a new game room
func CreateGameRoom(player1, player2 int) int {
	roomID := generateRoomID() // Generate a unique room ID
	ws.ActiveGames[roomID] = &ws.Game{
		RoomID:     roomID,
		WhiteID:    player1,
		BlackID:    player2,
		BoardState: InitialBoardState,
		Turn:       "white",
		Status:     "in-progress",
	}
	log.Println("LOOK AT ME")
	log.Printf("Created game room %d with players %d and %d", roomID, player1, player2)
	log.Printf("Active games: %v", ws.ActiveGames)
	return roomID
}

var uniqueIdCounter int = 0
var idMutex sync.Mutex

func generateRoomID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	uniqueIdCounter++
	return uniqueIdCounter

}
