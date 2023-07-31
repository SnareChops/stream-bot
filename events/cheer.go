package events

import (
	"encoding/json"
	"fmt"
)

const CheerAuthScope = "bits:read"

// done
type CheerEvent struct {
	IsAnon               bool   `json:"is_anonymous"`           // Whether the user cheered anonymously or not
	UserID               string `json:"user_id"`                // The user ID for the user now following the specified channel.
	UserLogin            string `json:"user_login"`             // The user login for the user now following the specified channel.
	UserName             string `json:"user_name"`              // The user display name for the user now following the specified channel.
	BroadcasterUserID    string `json:"broadcaster_user_id"`    // The requested broadcaster ID
	BroadcasterUserLogin string `json:"broadcaster_user_login"` // The requested broadcaster login
	BroadcasterUserName  string `json:"broadcaster_user_name"`  // The requested broadcaster display name
	Message              string `json:"message"`                // The message sent with the cheer
	Bits                 int    `json:"bits"`                   // The number of bits cheered
}

// done
type CheerSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"` // Empty string if anon
	Amount   int    `json:"amount"`
	Message  string `json:"message"`
}

// done
func NewCheerEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.cheer",
		Version: "1",
		Condition: map[string]interface{}{
			"broadcaster_user_id": userId,
		},
		Handler: handleCheer,
	}
}

// done
func handleCheer(event map[string]interface{}) (result []byte) {
	println("Handling cheer event")
	// Decode event
	cheer := &CheerEvent{}
	err := decode(event, cheer)
	if err != nil {
		fmt.Printf("Failed to decode CheerEvent: %s\n", err)
		return
	}
	// Create UI signal
	result, err = json.Marshal(CheerSignal{
		Kind:     "twitch.cheer",
		Username: cheer.UserName,
		Amount:   cheer.Bits,
		Message:  cheer.Message,
	})
	if err != nil {
		fmt.Printf("Failed to encode CheerSignal: %s\n", err)
		return
	}
	fmt.Printf("CheerSignal %s\n", result)
	return
}
