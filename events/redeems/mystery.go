package redeems

type MysterySignal struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func Mystery(signal RedeemSignal) []byte {
	return []byte{}
}
