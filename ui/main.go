package ui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
	"nhooyr.io/websocket"
)

func Start() {
	addui := make(chan *websocket.Conn)
	addadmin := make(chan *websocket.Conn)
	go uiemitter(addui)
	go adminemitter(addadmin)

	router := chi.NewRouter()

	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.html")
	}))

	router.Get("/index.js", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.js")
	}))

	router.Get("/index.css", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "ui/index.css")
	}))

	router.Get("/ws", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := websocket.Accept(w, req, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		addui <- ws
	}))

	router.Get("/adminws", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := websocket.Accept(w, req, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		addadmin <- ws
	}))

	router.Get("/files/{name}", func(w http.ResponseWriter, req *http.Request) {
		name, err := url.PathUnescape(chi.URLParam(req, "name"))
		if err != nil {
			fmt.Printf("Failed to unescape file name: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Serving file %s\n", name)
		http.ServeFile(w, req, "files/"+name)
	})

	router.Get("/video/{name}", func(w http.ResponseWriter, req *http.Request) {
		name, err := url.PathUnescape(chi.URLParam(req, "name"))
		if err != nil {
			fmt.Printf("Failed to unescape file name: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		video, err := os.Open("files/video/" + name)
		if err != nil {
			fmt.Printf("Failed to open video file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer video.Close()
		w.Header().Set("Content-Type", "video/mp4")
		fmt.Printf("Serving video %s\n", name)
		http.ServeContent(w, req, name, time.Now(), video)
	})

	log.Print(http.ListenAndServe(":3000", router))
}

func uiemitter(add chan *websocket.Conn) {
	connections := []*websocket.Conn{}
	for {
		select {
		case message := <-signals.SendToUI:
			fmt.Printf("Sending to UI: %s\n", string(message))
			for i, conn := range connections {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx, websocket.MessageText, message)
				if err != nil {
					fmt.Printf("Failed to send message to UI: %s\n", err)
					connections = slices.Delete(connections, i, i+1)
				}
				cancel()
			}
		case <-signals.Shutdown:
			for _, conn := range connections {
				conn.Close(websocket.StatusNormalClosure, "goodbye")
			}
			return
		case connection := <-add:
			println("UI Client connected!")
			connections = append(connections, connection)
		}
	}
}

func adminemitter(add chan *websocket.Conn) {
	connections := []*websocket.Conn{}
	for {
		select {
		case message := <-signals.SendToAdmin:
			fmt.Printf("Sending to Admin: %s\n", string(message))
			for _, conn := range connections {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx, websocket.MessageText, message)
				cancel()
				if err != nil {
					fmt.Printf("Failed to send message to admin: %s\n", err)
				}
			}
		case <-signals.Shutdown:
			for _, conn := range connections {
				conn.Close(websocket.StatusNormalClosure, "goodbye")
			}
			return
		case connection := <-add:
			connections = append(connections, connection)
		}
	}
}
