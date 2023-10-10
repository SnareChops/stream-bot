package irc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type SubSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Tier     int    `json:"tier"`
	Gift     bool   `json:"gift"`
	Streak   int    `json:"streak"`
}

func sub(message twitch.UserNoticeMessage) {
	streak, err := strconv.ParseInt(message.MsgParams["msg-param-streak-months"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse sub streak: %s\n", err)
		streak = 0
	}
	signal, err := json.Marshal(SubSignal{
		Kind:     "twitch.sub",
		Username: message.User.DisplayName,
		Tier:     transformTier(message.MsgParams["msg-param-sub-plan"]),
		Streak:   int(streak),
	})
	if err != nil {
		fmt.Printf("Failed to marshal SubSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}

func transformTier(tier string) int {
	switch tier {
	case "Prime":
		return 0
	case "1000":
		return 1
	case "2000":
		return 2
	case "3000":
		return 3
	default:
		return 0
	}
}
