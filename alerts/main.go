package alerts

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

func NewRouter() chi.Router {
	emitter := signals.NewSignalEmitter("Alerts", signals.SendToUI)
	go emitter.Start()

	router := chi.NewRouter()
	router.Get("/", file("alerts/index.html"))
	router.Get("/index.js", file("alerts/index.js"))
	router.Get("/index.css", file("alerts/index.css"))
	router.Get("/assets/{name}", func(w http.ResponseWriter, req *http.Request) {
		name, err := url.PathUnescape(chi.URLParam(req, "name"))
		if err != nil {
			fmt.Printf("Failed to unescape file name: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, req, "alerts/assets/"+name)
	})

	router.Get("/ws", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ws, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		emitter.Add(ws)
	}))

	return router
}

func file(path string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	})
}
