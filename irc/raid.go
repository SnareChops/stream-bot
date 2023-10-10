package irc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type RaidSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Viewers  int    `json:"viewers"`
}

func raid(message twitch.UserNoticeMessage) {
	viewers, err := strconv.ParseInt(message.MsgParams["msg-param-viewerCount"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse viewer count: %s\n", err.Error())
	}
	signal, err := json.Marshal(RaidSignal{
		Kind:     "twitch.raid",
		Username: message.MsgParams["msg-param-displayName"],
		Viewers:  int(viewers),
	})
	if err != nil {
		fmt.Printf("Failed to marshal RaidSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}
