package events

import (
	"encoding/json"
	"fmt"
)

// done
type GiftEvent struct {
	UserID               string `json:"user_id"`                // The user ID of the user who sent the subscription gift. Set to null if it was an anonymous subscription gift
	UserLogin            string `json:"user_login"`             // The user login of the user who sent the gift. Set to null if it was an anonymous subscription gift
	UserName             string `json:"user_name"`              // The user display name of the user who sent the gift. Set to null if it was an anonymous subscription gift
	BroadcasterUserID    string `json:"broadcaster_user_id"`    // The broadcaster user ID
	BroadcasterUserLogin string `json:"broadcaster_user_login"` // The broadcaster login
	BroadcasterUserName  string `json:"broadcaster_user_name"`  // The broadcaster display name
	Total                int    `json:"total"`                  // The number of subscriptions in the subscription gift
	Tier                 string `json:"tier"`                   // The tier of subscriptions in the subscription gift
	CumulativeTotal      int    `json:"cumulative_total"`       // The number of subscriptions gifted by this user in the channel. This value is null for anonymous gifts or if the gifter has opted out of sharing this information
	IsAnon               bool   `json:"is_anonymous"`           // Whether the subscription gift was anonymous
}

// done
type GiftSignal struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
	Tier     int    `json:"tier"`
	Amount   int    `json:"amount"`
	Total    int    `json:"total"`
	Anon     bool   `json:"anon"`
}

func NewGiftEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.subscription.gift",
		Version: "1",
		Condition: map[string]interface{}{
			"broadcaster_user_id": userId,
		},
		Handler: handleGift,
	}
}

func handleGift(event map[string]interface{}) (result []byte) {
	println("Handling gift event")
	// Decode gift
	gift := &GiftEvent{}
	err := decode(event, gift)
	if err != nil {
		fmt.Printf("Failed to decode GiftEvent: %s\n", err)
		return
	}
	// Encode UI GiftSignal
	result, err = json.Marshal(GiftSignal{
		Kind:     "twitch.gift",
		Username: gift.UserName,
		Tier:     transformTier(gift.Tier),
		Amount:   gift.Total,
		Total:    gift.CumulativeTotal,
		Anon:     gift.IsAnon,
	})
	if err != nil {
		fmt.Printf("Failed to encode GiftSignal: %s\n", err)
		return
	}
	fmt.Printf("GiftSignal %s\n", result)
	return
}
