package ws

import "encoding/json"

func decodePayload[T any](payload map[string]interface{}, out *T) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, out)
}
