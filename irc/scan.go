package irc

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SnareChops/twitchbot/signals"
	"golang.org/x/exp/slices"
)

type AudioSignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func scan(content []string) {
	if slices.Contains(content, "GOODONE") {
		audio("audio-emotes/GOODONE.mp3")
		return
	}
	if slices.Contains(content, "CLUELESS") {
		audio("audio-emotes/CLUELESS.mp3")
		return
	}
	if slices.Contains(content, "PEEK") {
		audio("audio-emotes/PEEK.mp3")
		return
	}
	if slices.Contains(content, "SEX") {
		audio("audio-emotes/SEX.mp3")
		return
	}
	if slices.Contains(content, "WOW") {
		audio("audio-emotes/WOW.mp3")
		return
	}
}

func audio(file string) {
	data, err := json.Marshal(&AudioSignal{
		Kind: "audio",
		URL:  "files/" + url.PathEscape(file),
	})
	if err != nil {
		fmt.Printf("Failed to marshal AudioSignal: %s\n", err)
		return
	}
	signals.SendToUI <- data
}
