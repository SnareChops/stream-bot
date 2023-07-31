package events

import (
	"encoding/json"
	"fmt"
)

const FollowAuthScope = "moderator:read:followers"

type FollowEvent struct {
	UserID               string `json:"user_id"`                // The user ID for the user now following the specified channel.
	UserLogin            string `json:"user_login"`             // The user login for the user now following the specified channel.
	UserName             string `json:"user_name"`              // The user display name for the user now following the specified channel.
	BroadcasterUserID    string `json:"broadcaster_user_id"`    // The requested broadcaster ID
	BroadcasterUserLogin string `json:"broadcaster_user_login"` // The requested broadcaster login
	BroadcasterUserName  string `json:"broadcaster_user_name"`  // The requested broadcaster display name
	FollowedAt           string `json:"followed_at"`            // RFC3399 timestamp of when the follow occurred.
}

type FollowSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
}

func NewFollowEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.follow",
		Version: "2",
		Condition: map[string]interface{}{
			"broadcaster_user_id": userId,
			"moderator_user_id":   userId,
		},
		Handler: handleFollow,
	}
}

func handleFollow(event map[string]interface{}) (result []byte) {
	println("Handling follow event")
	// Decode follow
	follow := FollowEvent{}
	err := decode(event, &follow)
	if err != nil {
		fmt.Printf("Failed to decode FollowEvent: %s\n", err)
		return
	}

	// Encode UI FollowSignal
	result, err = json.Marshal(FollowSignal{
		Kind:     "twitch.follow",
		Username: follow.UserName,
	})
	if err != nil {
		fmt.Printf("Failed to encode FollowSignal: %s\n", err)
		return
	}
	fmt.Printf("FollowSignal %s\n", result)
	return
}
