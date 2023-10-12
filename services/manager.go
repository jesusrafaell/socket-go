package services

import (
	"crashsaver/websocket/data"
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
		CheckOrigin:     checkOrigin,
	}
)

// all true
func checkOrigin(r *http.Request) bool {
	return true
}

type Manager struct {
	clients   ClientList
	incidents *data.IncidentList
	sync.RWMutex
}

func NewManager(incidents *data.IncidentList) *Manager {
	return &Manager{
		clients:   make(ClientList),
		incidents: incidents,
	}
}

func (m *Manager) ServerWS(c echo.Context) error {
	// Upgrade  HTTP to WebSocket
	ws, err := webSocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Error upgrading to Manager: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Create New Client and add client
	client := NewClient(ws, m)
	m.addClient(client)

	log.Printf("new connecton \n")

	//Start client process
	go client.readMessages(m)
	// go client.writeMessages()

	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		log.Printf("Manager, removeClient(): -> %s\n", client.conn.RemoteAddr().String())
		//close connection
		client.conn.Close()
		//reove client
		delete(m.clients, client)
	}
}

// send msg all clients
func (m *Manager) sendMessageToAllClients(message string) {
	m.Lock()
	defer m.Unlock()

	for clientws := range m.clients {
		err := clientws.writeMessage(message)
		// log.Printf("send msg (%s) to (%s)\n", message, clientws.conn.LocalAddr().String())
		if err != nil {
			log.Printf("Error sending message to client %s: %v", clientws.conn.LocalAddr().String(), err)
		}
	}
}

func (m *Manager) sendMessageToClient(client *Client, message string) {
	client.writeMessage(message)
}
