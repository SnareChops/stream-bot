package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

func connect(url string) (*websocket.Conn, string, error) {
	fmt.Printf("Connecting WebSocket to %s\n", url)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ws, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, "", err
	}
	_, data, err := ws.Read(context.Background())
	if err != nil {
		return nil, "", err
	}

	welcome := &Message{}
	err = json.Unmarshal(data, welcome)
	if err != nil {
		return nil, "", err
	}

	payload := &WelcomePayload{}
	err = decodePayload(welcome.Payload, payload)
	if err != nil {
		return nil, "", err
	}
	fmt.Printf("Connected with session ID: %s\n", payload.Session.ID)
	return ws, payload.Session.ID, nil
}
