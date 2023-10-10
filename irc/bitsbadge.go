package irc

import (
	"encoding/json"
	"fmt"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type BitsBadgeSignal struct {
	Kind     string  `json:"kind"`
	Username string  `json:"username"`
	Tier     int     `json:"tier"`
	Message  string  `json:"message"`
	Emotes   []Emote `json:"emotes"`
}

func bitsbadgetier(message twitch.UserNoticeMessage) {
	signal, err := json.Marshal(BitsBadgeSignal{
		Kind:     "twitch.bitsbadge",
		Username: message.User.DisplayName,
		Tier:     transformTier(message.MsgParams["msg-param-threshold"]),
		Message:  message.Message,
		Emotes:   transformEmotes(message.Emotes),
	})
	if err != nil {
		fmt.Printf("Failed to parse BitsBadgeSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}
