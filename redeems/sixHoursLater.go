package redeems

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type SixHoursLaterSignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func SixHoursLater(signal RedeemSignal) []byte {
	out, err := json.Marshal(&SixHoursLaterSignal{
		Kind: "audio",
		URL:  "files/" + url.PathEscape("audio/6-hrs-later.mp3"),
	})
	if err != nil {
		fmt.Printf("Failed to encode six hours later redeem: %s\n", err)
	}
	return out
}
