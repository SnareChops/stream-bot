package redeems

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type RedeemReward struct {
	ID     string `json:"id"`     // The reward identifier
	Title  string `json:"title"`  // The reward name
	Cost   int    `json:"cost"`   // The reward cost
	Prompt string `json:"prompt"` // The reward description
}

type RedeemSignal struct {
	Kind     string       `json:"kind"`
	Username string       `json:"username"`
	Input    string       `json:"input"`
	Reward   RedeemReward `json:"reward"`
}

func ProcessRedeem(signal RedeemSignal) []byte {
	fmt.Printf("Handle redeem %s\n", signal.Reward.Title)
	switch strings.ToLower(signal.Reward.Title) {
	case "hi":
		return Hi(signal)
	case "hooray":
		return Hooray(signal)
	case "bass solo":
		return BassSolo(signal)
	default:
		return nil
	}
}
