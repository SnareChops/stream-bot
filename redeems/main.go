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
	case "bass solo":
		return BassSolo(signal)
	case "hooray":
		return Audio("audio/hooray.mp3")
	case "six hours later":
		return Audio("audio/six-hours-later.mp3")
	case "amazing":
		return Audio("audio/amazing-how-wonderful.mp3")
	case "disgusting pervert":
		return Audio("audio/disgusting-pervert-absolutely-unacceptable.mp3")
	case "i'm bored":
		return Audio("audio/im-bored.mp3")
	case "shut up":
		return Audio("audio/omg-shut-up.mp3")
	case "meanwhile":
		return Audio("audio/spy-x-flute.mp3")
	case "stop it":
		return Audio("audio/stop-it.mp3")
	case "that won't work":
		return Audio("audio/that-wont-work-for-you.mp3")
	case "very impressive":
		return Audio("audio/very-impressive-indeed.mp3")
	case "dummy":
		return Audio("audio/what-a-complete-moron.mp3")
	default:
		return nil
	}
}
