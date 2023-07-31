package ws

import (
	"context"
	"encoding/json"
	"time"

	"nhooyr.io/websocket"
)

func connect(url string) (*websocket.Conn, string, error) {
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

	return ws, payload.Session.ID, nil
}
