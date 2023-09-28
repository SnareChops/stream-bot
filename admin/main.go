package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"nhooyr.io/websocket"
)

func Start() {
	SyncMysteries()
}

type AdminSignal struct {
	Kind string   `json:"kind"`
	Args []string `json:"args"`
}

func Listen(connection *websocket.Conn) {
	for {
		_, message, err := connection.Read(context.Background())
		if err != nil {
			fmt.Printf("Failed to read message from admin websocket: %s\n", err)
		}
		signal := new(AdminSignal)
		err = json.Unmarshal(message, signal)
		if err != nil {
			fmt.Printf("Failed to unmarshal admin signal json: %s\n", err)
		}
		switch signal.Kind {
		default:
			fmt.Printf("Unrecognized admin signal kind: %s\n", signal.Kind)
		}
	}
}
