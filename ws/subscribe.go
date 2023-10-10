package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SubscribeBody struct {
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Condition map[string]interface{} `json:"condition"`
	Transport SubscribeBodyTransport `json:"transport"`
}

type SubscribeBodyTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

type SubscribeResponse struct {
	Data         []SubscribeResponseData `json:"data"`
	Total        int                     `json:"total"`
	TotalCost    int                     `json:"total_cost"`
	MaxTotalCost int                     `json:"max_total_cost"`
}

type SubscribeResponseData struct {
	ID        string                     `json:"id"`
	Status    string                     `json:"status"`
	Type      string                     `json:"type"`
	Version   string                     `json:"version"`
	Condition map[string]interface{}     `json:"condition"`
	CreatedAt string                     `json:"created_at"`
	Transport SubscribeResponseTransport `json:"transport"`
	Cost      int                        `json:"cost"`
}

type SubscribeResponseTransport struct {
	Method      string `json:"method"`
	SessionID   string `json:"session_id"`
	ConnectedAt string `json:"connected_at"`
}

func subscribe(clientId, token, subscribeUrl string, body SubscribeBody) (*SubscribeResponse, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	readerBody := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, subscribeUrl, readerBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", clientId)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %v", res.Status, res.Body)
	}

	resBody := SubscribeResponse{}
	err = json.NewDecoder(res.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &resBody, nil
}
