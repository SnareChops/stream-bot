package events

import (
	"fmt"

	"github.com/SnareChops/twitchbot/events/redeems"
)

const RedeemAuthScope = "channel:read:redemptions"

// done
type RedeemEvent struct {
	ID                   string               `json:"id"`                     // The redemption identifier
	BroadcasterUserID    string               `json:"broadcaster_user_id"`    // The requested broadcaster ID
	BroadcasterUserLogin string               `json:"broadcaster_user_login"` // The requested broadcaster login
	BroadcasterUserName  string               `json:"broadcaster_user_name"`  // The requested broadcaster display name
	UserID               string               `json:"user_id"`                // The user ID for the user now following the specified channel.
	UserLogin            string               `json:"user_login"`             // The user login for the user now following the specified channel.
	UserName             string               `json:"user_name"`              // The user display name for the user now following the specified channel.
	UserInput            string               `json:"user_input"`             // The user input provided. Empty string if not provided
	Status               string               `json:"status"`                 // Defaults to "unfulfilled"
	Reward               redeems.RedeemReward `json:"reward"`                 // Basic information about the reward that was redeemed at the time it was redeemed
	RedeemedAt           string               `json:"reedeemed_at"`           // RFC3399 timestamp of when the reward was redeemed
}

func NewRedeemEventSub(userId string) *EventSub {
	return &EventSub{
		Type:    "channel.channel_points_custom_reward_redemption.add",
		Version: "1",
		Condition: map[string]interface{}{
			"broadcaster_user_id": userId,
		},
		Handler: handleRedeem,
	}
}

// done
func handleRedeem(event map[string]interface{}) []byte {
	println("Handling redeem event")
	// Decode redeem
	redeem := &RedeemEvent{}
	err := decode(event, redeem)
	if err != nil {
		fmt.Printf("Failed to decode RedeemEvent: %s\n", err)
		return nil
	}
	// Route redeem
	return redeems.ProcessRedeem(redeems.RedeemSignal{
		Kind:     "twitch.redeem",
		Username: redeem.UserName,
		Input:    redeem.UserInput,
		Reward:   redeem.Reward,
	})
}
