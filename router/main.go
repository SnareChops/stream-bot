package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/SnareChops/twitchbot/admin"
	"github.com/SnareChops/twitchbot/alerts"
	"github.com/SnareChops/twitchbot/avatar"
	"github.com/SnareChops/twitchbot/ui"
	"github.com/go-chi/chi/v5"
)

func Start() {

	router := chi.NewRouter()

	router.Mount("/admin", admin.NewRouter())
	router.Mount("/alerts", alerts.NewRouter())
	router.Mount("/avatar", avatar.NewRouter())
	router.Mount("/ui", ui.NewRouter())

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
