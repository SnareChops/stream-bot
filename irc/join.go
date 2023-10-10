package irc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type JoinSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Time     int64  `json:"time"`
}

func join(message twitch.UserJoinMessage) {
	signal, err := json.Marshal(JoinSignal{
		Kind:     "twitch.join",
		Username: message.User,
		Time:     time.Now().UnixMilli(),
	})
	if err != nil {
		fmt.Printf("Failed to create join signal: %s\n", err.Error())
		return
	}
	signals.SendToAdmin <- signal
}
