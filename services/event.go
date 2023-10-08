package services

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

// const (
// 	// EventSendMessage is the event name for new chat messages sent
// 	EventSendMessage = "send_message"
// 	// EventNewMessage is a response to send_message
// 	EventNewMessage = "new_message"
// 	// EventChangeRoom is event when switching rooms
// 	EventChangeRoom = "change_room"
// )

// type SendMessageEvent struct {
// 	Message string `json:"message"`
// 	From    string `json:"from"`
// }

// type NewMessageEvent struct {
// 	SendMessageEvent
// 	Sent time.Time `json:"sent"`
// }

// func SendMessageHandler(event Event, c *Client) error {
// 	// Marshal Payload into wanted format
// 	var chatevent SendMessageEvent
// 	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
// 		return fmt.Errorf("bad payload in request: %v", err)
// 	}

// 	var broadMessage NewMessageEvent

// 	broadMessage.Sent = time.Now()
// 	broadMessage.Message = chatevent.Message
// 	broadMessage.From = chatevent.From

// 	data, err := json.Marshal(broadMessage)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal broadcast message: %v", err)
// 	}

// 	// Place payload into an Event
// 	var outgoingEvent Event
// 	outgoingEvent.Payload = data
// 	outgoingEvent.Type = EventNewMessage
// 	// Broadcast to all other Clients
// 	for client := range c.manager.clients {
// 		// Only send to clients inside the same chatroom
// 		if client.chatroom == c.chatroom {
// 			client.egress <- outgoingEvent
// 		}

// 	}
// 	return nil
// }

// type ChangeRoomEvent struct {
// 	Name string `json:"name"`
// }

// // ChatRoomHandler will handle switching of chatrooms between clients
// func ChatRoomHandler(event Event, c *Client) error {
// 	// Marshal Payload into wanted format
// 	var changeRoomEvent ChangeRoomEvent
// 	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
// 		return fmt.Errorf("bad payload in request: %v", err)
// 	}

// 	// Add Client to chat room
// 	c.chatroom = changeRoomEvent.Name

// 	return nil
// }
