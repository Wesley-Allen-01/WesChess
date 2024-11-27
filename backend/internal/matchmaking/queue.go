package matchmaking

import (
	"sync"
	"log"
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
	log.Println("Attempting to lock queueMutex")
    queueMutex.Lock()
    defer func() {
        log.Println("Unlocking queueMutex")
        queueMutex.Unlock()
    }()

	if len(matchQueue) < 2 {
		return 0, 0, 0, false
	}

	player1 := matchQueue[0]
	player2 := matchQueue[1]
	matchQueue = matchQueue[2:]

	roomID := CreateGameRoom(player1, player2)

	MatchedMutex.Lock()
	defer MatchedMutex.Unlock()
	MatchedPlayers[player1] = roomID
	MatchedPlayers[player2] = roomID

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
