package matchmaking

import (
	"sync"
)

var matchQueue []int
var queueMutex sync.Mutex

func EnqueuePlayer(playerID int) {
	queueMutex.Lock()
	defer queueMutex.Unlock()
	matchQueue = append(matchQueue, playerID)
}


var MatchedPlayers = make(map[int]int) // map from userID to roomID
var MatchedMutex sync.Mutex

func MatchPlayers() (int, int, int, bool) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(matchQueue) < 2 {
		return 0, 0, 0, false
	}

	player1 := matchQueue[0]
	player2 := matchQueue[1]
	matchQueue = matchQueue[2:]

	roomID := generateRoomID()

	MatchedMutex.Lock()
	defer MatchedMutex.Unlock()
	MatchedPlayers[player1] = roomID
	MatchedPlayers[player2] = roomID
	MatchedMutex.Unlock()

	return player1, player2, roomID, true
}

func CheckUserMatch(userID int) (int, bool) {
	MatchedMutex.Lock()
	defer MatchedMutex.Unlock()
	roomID, exists := MatchedPlayers[userID]
	if !exists {
		return 0, false
	}
	return roomID, true
}
