package events

import (
	"encoding/json"
	"fmt"
)

type ResubEvent struct {
	UserID               string       `json:"user_id"`                // The user ID of the user who sent a resubscription chat message
	UserLogin            string       `json:"user_login"`             // The user login of the user who sent a resubscription chat message
	UserName             string       `json:"user_name"`              // The user display name of the user who sent a resubscription chat message
	BroadcasterUserID    string       `json:"broadcaster_user_id"`    // The broadcaster user ID
	BroadcasterUserLogin string       `json:"broadcaster_user_login"` // The broadcaster login
	BroadcasterUserName  string       `json:"broadcaster_user_name"`  // The broadcaster display name
	Tier                 string       `json:"tier"`                   // The tier of the user's subscription (1000, 2000, 3000)
	Message              ResubMessage `json:"message"`                // An object that contains the resubscription message and emote information needed to recreate the message
	CumulativeMonts      int          `json:"cumulative_monts"`       // The total number of months the user has been subscribed to the channel
	StreakMonths         int          `json:"streak_months"`          // The number of consecutive months the user's current subscription has been active. This value is nul if the user has opted out of sharing this information
	DurationMonths       int          `json:"duration_months"`        // The month duration of the subscription
}

// done
type ResubMessage struct {
	Text   string  `json:"text"`   // The text of the resubscription chat message
	Emotes []Emote `json:"emotes"` // An array that includes the emote ID and start and end positions for where the emote appears in the text
}

// done
type Emote struct {
	Begin int    `json:"begin"` // The index of where the Emote starts in the text
	End   int    `json:"end"`   // The index of where the Emote ends in the text
	ID    string `json:"id"`    // The emote ID
}

// done
type ResubSignal struct {
	Kind     string       `json:"kind"`
	Username string       `json:"username"`
	Tier     int          `json:"tier"`
	Months   int          `json:"months"`
	Message  ResubMessage `json:"message"`
}

func NewResubEventSub(userId string) *EventSub {
	return &EventSub{
		Type:      "channel.subscription.message",
		Version:   "1",
		Condition: map[string]interface{}{},
		Handler:   handleResub,
	}
}

func handleResub(event map[string]interface{}) (result []byte) {
	println("Handling resub event")
	// Decide resub
	resub := &ResubEvent{}
	err := decode(event, resub)
	if err != nil {
		fmt.Printf("Failed to decode ResubEvent: %s\n", err)
		return
	}
	// Encode UI ResubSignal
	result, err = json.Marshal(ResubSignal{
		Kind:     "twitch.resub",
		Username: resub.UserName,
		Tier:     transformTier(resub.Tier),
		Months:   resub.CumulativeMonts,
		Message:  resub.Message,
	})
	if err != nil {
		fmt.Printf("Failed to encode ResubSignal: %s\n", err)
		return
	}
	fmt.Printf("ResubSignal %s\n", result)
	return
}
