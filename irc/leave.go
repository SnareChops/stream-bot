package irc

import (
	"encoding/json"
	"fmt"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type LeaveSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
}

func leave(message twitch.UserPartMessage) {
	signal, err := json.Marshal(LeaveSignal{
		Kind:     "twitch.leave",
		Username: message.User,
	})
	if err != nil {
		fmt.Printf("Failed to create leave signal: %s\n", err.Error())
		return
	}
	signals.SendToAdmin <- signal
}
