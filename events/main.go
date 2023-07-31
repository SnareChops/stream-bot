package events

import "encoding/json"

const (
	SubscriptionTypeChannelFollow = "channel.follow"
)

type EventSub struct {
	Type      string
	Version   string
	Condition map[string]interface{}
	Handler   func(event map[string]interface{}) []byte
}

func decode[T any](payload map[string]interface{}, out T) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, out)
}
