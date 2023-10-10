package irc

import (
	"encoding/json"
	"fmt"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type RitualSignal struct {
	Kind     string  `json:"kind"`
	Type     string  `json:"type"`
	Username string  `json:"username"`
	Message  string  `json:"message"`
	Emotes   []Emote `json:"emotes"`
}

func ritual(message twitch.UserNoticeMessage) {
	signal, err := json.Marshal(RitualSignal{
		Kind:     "twitch.ritual",
		Type:     message.MsgParams["msg-param-ritual-name"],
		Username: message.User.DisplayName,
		Message:  message.Message,
		Emotes:   transformEmotes(message.Emotes),
	})
	if err != nil {
		fmt.Printf("Failed to marshal RitualSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}
