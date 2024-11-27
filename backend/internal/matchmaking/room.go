package matchmaking

import (
	"WesChess/backend/internal/ws"
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