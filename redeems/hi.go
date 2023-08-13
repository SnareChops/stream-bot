package redeems

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

var HiFiles []string

func init() {
	files, err := os.ReadDir("files/hi")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		HiFiles = append(HiFiles, file.Name())
	}
}

type HiSignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func Hi(signal RedeemSignal) []byte {
	choice := random.Intn(len(HiFiles))
	out, err := json.Marshal(&HiSignal{
		Kind: "audio",
		URL:  "files/" + url.PathEscape("hi/"+HiFiles[choice]),
	})
	if err != nil {
		fmt.Printf("Failed to encode hi redeem: %s\n", err)
	}
	return out
}
