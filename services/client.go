package services

import (
	"crashsaver/websocket/data"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

var NextID int

type Client struct {
	id      int
	conn    *websocket.Conn
	manager *Manager

	//egress is used to concurrecnt writes on the webscoket connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	NextID++
	return &Client{
		id:      NextID,
		conn:    conn,
		manager: manager,
		egress:  make(chan []byte),
	}
}

func (c *Client) readMessages(m *Manager) {
	defer func() {
		//cleanup connecton
		c.manager.removeClient(c)
	}()

	for {

		var payload data.WebSocketMessage
		err := c.conn.ReadJSON(&payload)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			} else {
				log.Printf("Client, readMessages(): %v", err)
			}
			break
		}

		// log.Printf("Recibed from client: %s\n", string(payload.Type))
		switch payload.Type {
		case "create", "update", "delete", "get":
			// Crear una instancia del servicio de Incidentes y manejar el mensaje
			incidentServices := NewIncidents(m)
			incidentServices.HandleWebSocketMessage(c, payload)
		default:
			log.Printf("Type not found readMessages(): %s\n", payload.Type)
		}
		////////////////////////////////////////

	}
}

func (c *Client) writeMessages() {
	// pongWait is how long we will await a pong response from client
	pongWait := 10 * time.Second
	pingInterval := (pongWait * 9) / 10
	ticker := time.NewTicker(pingInterval)

	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				// Manager has closed this connection channel,
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed: %v\n", err)
				}
				return
			}
			// Write a Regular text message to the connection
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v\n", err)
			}
			log.Println("message sent")
		case <-ticker.C:
			// log.Println("ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeMessage(message string) error {
	// log.Printf("send msg client %s\n", message)
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
