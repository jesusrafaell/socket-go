package socketService

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connecton     *websocket.Conn
	socketManager *SocketManager
}

func NewClient(conn *websocket.Conn, socketManager *SocketManager) *Client {
	return &Client{
		connecton:     conn,
		socketManager: socketManager,
	}
}

type WebSocketIncidentMovil struct {
	Msg string `json:"msg"`
}

func (c *Client) readMessages() {
	defer func() {
		//cleanup connecton
		c.socketManager.removeClient(c)
	}()

	fmt.Println("Reading messages...")

	for {

		messageType, payload, err := c.connecton.ReadMessage()

		fmt.Printf("User remote: -> %s\n", c.connecton.RemoteAddr().String())

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			} else {
				log.Printf("Client, readMessages(): %v", err)
			}
			break
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}
