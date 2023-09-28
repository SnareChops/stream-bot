package admin

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

var AdminTemplate *template.Template

func init() {
	var err error
	AdminTemplate, err = template.New("AdminTemplate").ParseFiles("admin/index.html")
	if err != nil {
		panic(err)
	}
}

func NewRouter() chi.Router {
	addadmin := make(chan *websocket.Conn)
	go adminemitter(addadmin)

	router := chi.NewRouter()

	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := AdminTemplate.ExecuteTemplate(w, "index", nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
		}
	}))
	router.Get("/index.js", file("admin/index.js"))
	router.Get("/index.css", file("admin/index.css"))
	// router.Get("/mysteries", GETMysteries)
	router.Get("/ws", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := websocket.Accept(w, req, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		addadmin <- ws
	}))
	return router
}

func file(path string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	})
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
			go Listen(connection)
		}
	}
}
