package ui

import (
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func Start(send chan []byte, close chan bool) {
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
		c, err := websocket.Accept(w, req, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		defer c.Close(websocket.StatusInternalError, "Unexpected closure")

		for {
			println("Waiting for send")
			select {
			case message := <-send:
				println("Received send")
				c.Write(req.Context(), websocket.MessageText, message)
			case <-close:
				c.Close(websocket.StatusNormalClosure, "goodbye")
			}
		}
	}))

	log.Print(http.ListenAndServe(":3000", nil))
}
