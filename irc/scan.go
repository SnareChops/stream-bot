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
	if slices.Contains(content, "80085") {
		signals.SendToIRC <- "nyaoki"
	}
	if slices.Contains(content, "Clap") {
		audio("audio-emotes/CLAP.mp3")
		return
	}
	if slices.Contains(content, "CLUELESS") {
		audio("audio-emotes/CLUELESS.mp3")
		return
	}
	if slices.Contains(content, "DIESOFCRINGE") {
		audio("audio-emotes/DIESOFCRINGE.mp3")
		return
	}
	if slices.Contains(content, "FLASHBANG") {
		audio("audio-emotes/FLASHBANG.mp3")
		return
	}
	if slices.Contains(content, "GOODONE") {
		audio("audio-emotes/GOODONE.mp3")
		return
	}
	if slices.Contains(content, "guraBONGO") {
		audio("audio-emotes/guraBONGO.mp3")
		return
	}
	if slices.Contains(content, "guraSing") {
		audio("audio-emotes/guraSing.mp3")
		return
	}
	if slices.Contains(content, "HUH") {
		audio("audio-emotes/HUH.mp3")
		return
	}
	if slices.Contains(content, "NODDERS") {
		audio("audio-emotes/NODDERS.mp3")
		return
	}
	if slices.Contains(content, "NOOOO") {
		audio("audio-emotes/NOOOO.mp3")
		return
	}
	if slices.Contains(content, "NOPERS") {
		audio("audio-emotes/NOPERS.mp3")
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
