package irc

import (
	"encoding/json"
	"fmt"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type PrimeUpgradeSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
}

func primepaidupgrade(message twitch.UserNoticeMessage) {
	signal, err := json.Marshal(PrimeUpgradeSignal{
		Kind:     "twitch.upgrade.prime",
		Username: message.User.DisplayName,
	})
	if err != nil {
		fmt.Printf("Failed to marshal PrimeUpgradeSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}

type GiftUpgradeSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Gifter   string `json:"gifter"`
}

func giftpaidupgrade(message twitch.UserNoticeMessage) {
	signal, err := json.Marshal(GiftUpgradeSignal{
		Kind:     "twitch.upgrade.gift",
		Username: message.User.DisplayName,
		Gifter:   message.MsgParams["msg-param-sender-name"],
	})
	if err != nil {
		fmt.Printf("Failed to marshal GiftUpgradeSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}

type AnonGiftUpgradeSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
}

func anongiftpaidupgrade(message twitch.UserNoticeMessage) {
	signal, err := json.Marshal(AnonGiftUpgradeSignal{
		Kind:     "twitch.upgrade.anon",
		Username: message.User.DisplayName,
	})
	if err != nil {
		fmt.Printf("Failed to marshal AnonGiftUpgradeSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}
