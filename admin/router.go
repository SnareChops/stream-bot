package admin

import (
	"fmt"
	"net/http"
	"text/template"

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
	emitter := signals.NewSignalEmitter("Admin", signals.SendToAdmin)
	go emitter.Start()

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
		emitter.Add(ws)
	}))
	return router
}

func file(path string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	})
}
