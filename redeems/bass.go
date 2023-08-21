package redeems

import (
	"encoding/json"
	"fmt"
)

type BassSignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func BassSolo(signal RedeemSignal) []byte {
	out, err := json.Marshal(&BassSignal{
		Kind: "video.small",
		URL:  "video/BassSolo.mp4",
	})
	if err != nil {
		fmt.Printf("Failed to encode bass redeem: %s\n", err)
	}
	return out
}
