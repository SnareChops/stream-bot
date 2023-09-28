package utils

import "encoding/json"

func decode[T any](payload map[string]interface{}, out *T) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, out)
}
