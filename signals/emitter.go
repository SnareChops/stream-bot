package signals

import (
	"context"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

var lastId = 0

type SignalEmitter struct {
	name        string
	connections map[int]*websocket.Conn
	signal      chan []byte
}

func NewSignalEmitter(name string, signal chan []byte) *SignalEmitter {
	return &SignalEmitter{
		name:        name,
		connections: map[int]*websocket.Conn{},
		signal:      signal,
	}
}

func (self *SignalEmitter) Add(connection *websocket.Conn) {
	lastId += 1
	id := lastId
	self.connections[id] = connection
	fmt.Printf("%s connection added! (%d)\n", self.name, id)
}

func (self *SignalEmitter) Start() {
	for {
		select {
		case <-Shutdown:
			for _, conn := range self.connections {
				conn.Close(websocket.StatusNormalClosure, "goodbye")
			}
			return
		case message := <-self.signal:
			for id, conn := range self.connections {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx, websocket.MessageText, message)
				if err != nil {
					fmt.Printf("Failed to emit for %s (%d): %s\n", self.name, id, err.Error())
					conn.Close(websocket.StatusAbnormalClosure, "failure to send")
					delete(self.connections, id)
				}
				cancel()
			}
		}
	}
}
