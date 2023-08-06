package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SnareChops/twitchbot/events"
	"nhooyr.io/websocket"
)

type ReconnectPayload struct {
	Session SessionPayload `json:"session"`
}

type EventPayload struct {
	Event map[string]interface{} `json:"event"`
}

func listen(ws *websocket.Conn, send chan []byte, reconnect chan string, close chan bool, subscriptions []*events.EventSub) {
	println("Starting WebSocket listener...")
	for {
		// Check if the connection has been closed, if so stop listening
		select {
		case <-close:
			println("Lister closed")
			return
		default:
		}
		// Read message from socket
		_, data, err := ws.Read(context.Background())
		if err != nil {
			fmt.Printf("Failed to read websocket message: %s\n", err)
			continue
		}
		// Decode message
		message := &Message{}
		err = json.Unmarshal(data, message)
		if err != nil {
			fmt.Printf("Failed to unmarshal message: %s\n", err)
			continue
		}
		// Handle reconnect message
		if message.Metadata.MessageType == MessageTypeReconnect {
			payload := &ReconnectPayload{}
			err = decodePayload(message.Payload, payload)
			if err != nil {
				fmt.Printf("Failed to decode reconnect message payload: %s\n", err)
			}
			reconnect <- payload.Session.ReconnectUrl
			continue
		}
		// Handle notification message
		if message.Metadata.MessageType == MessageTypeNotification {
			for _, subscription := range subscriptions {
				if subscription.Type == message.Metadata.SubscriptionType {
					payload := &EventPayload{}
					err := decodePayload(message.Payload, payload)
					if err != nil {
						fmt.Printf("Failed to decode message payload for event type %s: %s\n", subscription.Type, err)
					}
					select {
					case send <- subscription.Handler(payload.Event):
					default:
					}
					break
				}
			}
		}
	}
}
