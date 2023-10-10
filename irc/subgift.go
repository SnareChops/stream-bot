package irc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type GiftSignal struct {
	Kind       string `json:"kind"`
	Username   string `json:"username"`
	Tier       int    `json:"tier"`
	Recipient  string `json:"recipient"`
	Months     int    `json:"months"`
	GiftMonths int    `json:"giftMonths"`
}

func subgift(message twitch.UserNoticeMessage) {
	months, err := strconv.ParseInt(message.MsgParams["msg-param-months"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse subgift months: %s\n", err.Error())
	}
	giftMonths, err := strconv.ParseInt(message.MsgParams["msg-param-gift-months"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse subgift giftMonths: %s\n", err.Error())
	}
	signal, err := json.Marshal(GiftSignal{
		Kind:       "twitch.gift",
		Username:   message.User.DisplayName,
		Tier:       transformTier(message.MsgParams["msg-param-sub-plan"]),
		Recipient:  message.MsgParams["msg-param-recipient-display-name"],
		Months:     int(months),
		GiftMonths: int(giftMonths),
	})
	if err != nil {
		fmt.Printf("Failed to marshal GiftSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}

type CommunityGiftSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Tier     int    `json:"tier"`
	Amount   int    `json:"amount"`
}

func submysterygift(message twitch.UserNoticeMessage) {
	amount, err := strconv.ParseInt(message.MsgParams["msg-param-mass-gift-count"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse mass gift count: %s\n", err.Error())
	}
	signal, err := json.Marshal(CommunityGiftSignal{
		Kind:     "twitch.gift.community",
		Username: message.User.DisplayName,
		Tier:     transformTier(message.MsgParams["msg-param-sub-plan"]),
		Amount:   int(amount),
	})
	if err != nil {
		fmt.Printf("Failed to marshal community gift signal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}
