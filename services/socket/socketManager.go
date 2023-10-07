package socketService

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin2,
	}
)

// checkOrigin will check origin and return true if its allowed
func checkOrigin2(r *http.Request) bool {
	//all true
	return true
}

// func checkOrigin(r *http.Request) bool {

// 	// Grab the request origin
// 	origin := r.Header.Get("Origin")

// 	switch origin {
// 	case "http://localhost:8080":
// 		return true
// 	default:
// 		return false
// 	}
// }

type SocketManager struct {
	clients ClientList
	sync.RWMutex
}

func NewSocket() *SocketManager {
	return &SocketManager{
		clients: make(ClientList),
	}
}

func (w *SocketManager) ServerWS(c echo.Context) error {
	ws, err := webSocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading to SocketManager: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	log.Printf("new connection -> %s\n", ws.RemoteAddr().String())

	defer ws.Close()

	client := NewClient(ws, w)

	w.addClient(client)

	length := len(w.clients)
	log.Printf("Connections : %d \n", length)

	//Start client process
	go client.readMessages()

	return nil
}

func (w *SocketManager) addClient(client *Client) {
	w.Lock()
	defer w.Unlock()

	w.clients[client] = true
}

func (w *SocketManager) removeClient(client *Client) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.clients[client]; ok {
		log.Printf("Manager, removeClient(): -> %s\n", client.connecton.RemoteAddr().String())
		client.connecton.Close()
		delete(w.clients, client)
	}
}
