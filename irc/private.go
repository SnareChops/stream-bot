package irc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type MessageSignal struct {
	Kind     string  `json:"kind"`
	Username string  `json:"username"`
	Message  string  `json:"message"`
	Emotes   []Emote `json:"emotes"`
}

func private(message twitch.PrivateMessage) {
	println("Received message")
	content := strings.Split(strings.TrimSpace(message.Message), " ")
	if len(content) == 0 {
		return
	}
	if content[0][0] != '!' {
		scan(content)
	} else {
		command, err := FindCommand(strings.ToLower(content[0]))
		if err != nil {
			fmt.Printf("Failed to find command: %s, %s", content[0], err)
		}
		if canRun(command.Permission, message) {
			signals.SendToIRC <- command.Response
		}
	}
	signal, err := json.Marshal(MessageSignal{
		Kind:     "twitch.message",
		Username: message.User.DisplayName,
		Message:  message.Message,
		Emotes:   transformEmotes(message.Emotes),
	})
	if err != nil {
		fmt.Printf("Failed to marshal MessageSignal: %s\n", err.Error())
		return
	}
	signals.SendToAdmin <- signal
}
