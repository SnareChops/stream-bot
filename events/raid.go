package events

import (
	"encoding/json"
	"fmt"
)

// done
type RaidEvent struct {
	FromUserID    string `json:"from_broadcaster_user_id"`    // The broadcaster ID that created the raid
	FromUserLogin string `json:"from_broadcaster_user_login"` // The broadcaster login that created the raid
	FromUserName  string `json:"from_broadcaster_user_name"`  // The broadcaster display name that created the raid
	ToUserID      string `json:"to_broadcaster_user_id"`      // The broadcaster ID that received the raid
	ToUserLogin   string `json:"to_broadcaster_user_login"`   // The broadcaster login that recieved the raid
	ToUserName    string `json:"to_broadcaster_user_name"`    // The broadcaster display name that received the raid
	Viewers       int    `json:"viewers"`                     // The number of viewers in the raid
}

// done
type RaidSignal struct {
	Kind    string `json:"kind"`
	From    string `json:"from"`
	Viewers int    `json:"viewers"`
}

func NewRaidEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.raid",
		Version: "1",
		Condition: map[string]interface{}{
			"to_broadcaster_user_id": userId,
		},
		Handler: handleRaid,
	}
}

// done
func handleRaid(event map[string]interface{}) (result []byte) {
	println("Handling raid event")
	// Decode raid
	raid := &RaidEvent{}
	err := decode(event, &raid)
	if err != nil {
		fmt.Printf("Failed to decode RaidEvent: %s\n", err)
		return
	}
	// Encode UI RaidSignal
	result, err = json.Marshal(RaidSignal{
		Kind:    "twitch.raid",
		From:    raid.FromUserName,
		Viewers: raid.Viewers,
	})
	if err != nil {
		fmt.Printf("Failed to encode RaidSignal: %s\n", err)
		return
	}
	fmt.Printf("RaidSignal %s\n", result)
	return
}
