package redeems

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type HooraySignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func Hooray(signal RedeemSignal) []byte {
	out, err := json.Marshal(&HooraySignal{
		Kind: "audio",
		URL:  "files/" + url.PathEscape("audio/hooray.mp3"),
	})
	if err != nil {
		fmt.Printf("Failed to encode hooray redeem: %s\n", err)
	}
	return out
}
