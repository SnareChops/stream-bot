package ui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

func Start(send chan []byte, close chan bool) {
	add := make(chan *websocket.Conn)
	go emitter(send, close, add)

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.html")
	}))

	http.Handle("/index.js", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.js")
	}))

	http.Handle("/index.css", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.css")
	}))

	http.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := websocket.Accept(w, req, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		add <- ws
	}))

	log.Print(http.ListenAndServe(":3000", nil))
}

func emitter(send chan []byte, close chan bool, add chan *websocket.Conn) {
	connections := []*websocket.Conn{}
	for {
		select {
		case message := <-send:
			fmt.Printf("Sending to UI: %s\n", string(message))
			for _, conn := range connections {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx, websocket.MessageText, message)
				if err != nil {
					fmt.Printf("Failed to send message")
				}
				cancel()
			}
		case <-close:
			for _, conn := range connections {
				conn.Close(websocket.StatusNormalClosure, "goodbye")
			}
			return
		case connection := <-add:
			connections = append(connections, connection)
		}
	}
}
