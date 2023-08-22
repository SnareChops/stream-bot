package redeems

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type AudioSignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func Audio(path string) []byte {
	out, err := json.Marshal(&AudioSignal{
		Kind: "audio",
		URL:  "files/" + url.PathEscape(path),
	})
	if err != nil {
		fmt.Printf("Failed to encode audio redeem: %s\n", err)
	}
	return out
}
