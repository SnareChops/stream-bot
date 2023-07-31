package events

import (
	"encoding/json"
	"fmt"
)

const SubscribeAuthScope = "channel:read:subscriptions"

// done
type SubEvent struct {
	UserID               string `json:"user_id"`                // The user ID for the user now following the specified channel.
	UserLogin            string `json:"user_login"`             // The user login for the user now following the specified channel.
	UserName             string `json:"user_name"`              // The user display name for the user now following the specified channel.
	BroadcasterUserID    string `json:"broadcaster_user_id"`    // The requested broadcaster ID
	BroadcasterUserLogin string `json:"broadcaster_user_login"` // The requested broadcaster login
	BroadcasterUserName  string `json:"broadcaster_user_name"`  // The requested broadcaster display name
	Tier                 string `json:"tier"`                   // The tier of the subscription. Valid values are 1000, 2000, and 3000
	IsGift               bool   `json:"is_gift"`                // Whether the subscription is a gift
}

// done
type SubSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Tier     int    `json:"tier"`
	Gift     bool   `json:"gift"`
}

func NewSubEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.subscribe",
		Version: "1",
		Condition: map[string]interface{}{
			"broadcaster_user_id": userId,
		},
		Handler: handleSub,
	}
}

func handleSub(event map[string]interface{}) (result []byte) {
	println("Handling SubscribeEvent")
	// Decode SubscribeEvent
	subscribe := &SubEvent{}
	err := decode(event, subscribe)
	if err != nil {
		fmt.Printf("Failed to decode SubscribeEvent: %s\n", err)
		return
	}
	// Encode SubscribeSignal
	result, err = json.Marshal(SubSignal{
		Kind:     "twitch.sub",
		Username: subscribe.UserName,
		Tier:     transformTier(subscribe.Tier),
		Gift:     subscribe.IsGift,
	})
	if err != nil {
		fmt.Printf("Failed to encode SubscribeSignal: %s\n", err)
		return
	}
	fmt.Printf("SubscribeSignal %s\n", result)
	return
}

func transformTier(tier string) int {
	switch tier {
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
