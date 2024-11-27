package ws

import (
	"log"
	"sync"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn *websocket.Conn
	Room string
}

type Room struct {
	Clients map[*Client]bool
	Mutex   sync.Mutex
}

type Game struct {
	RoomID	 int
	WhiteID  int
	BlackID  int
	BoardState string
	Turn     string
	Status   string
	winner   string
}

var ActiveGames = make(map[int]*Game)
var GameMutex sync.Mutex

var rooms = make(map[string]*Room)

func HandleConnection(w http.ResponseWriter, r *http.Request, roomID string) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }

    log.Printf("Successfully upgraded connection for room %s", roomID)

    client := &Client{Conn: conn, Room: roomID}
    addClientToRoom(client, roomID)

    defer func() {
        removeClientFromRoom(client, roomID)
        conn.Close()
    }()

    // Directly call handleMessages without a goroutine
    handleMessages(client)
}

func handleMessages(client *Client) {
    log.Println("handleMessages started")
    for {
        _, message, err := client.Conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message from client in room %s: %v", client.Room, err)
            break
        }
        log.Printf("Message received from client in room %s: %s", client.Room, string(message))
        broadcastToRoom(client.Room, message)
    }
    log.Println("handleMessages ended")
}

func addClientToRoom(client *Client, roomID string) {
	room, exists := rooms[roomID]
	if !exists {
		room = &Room{Clients: make(map[*Client]bool)}
		rooms[roomID] = room
	}
	room.Mutex.Lock()
	room.Clients[client] = true
	room.Mutex.Unlock()
}

func removeClientFromRoom(client *Client, roomID string) {
	room, exists := rooms[roomID]
	if exists {
		room.Mutex.Lock()
		delete(room.Clients, client)
		room.Mutex.Unlock()
	}
}

func broadcastToRoom(roomID string, message []byte) {
	room, exists := rooms[roomID]
	if exists {
		room.Mutex.Lock()
		for client := range room.Clients {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error writing message: ", err)
				client.Conn.Close()
				delete(room.Clients, client)
			}
		}
		room.Mutex.Unlock()
	}
}