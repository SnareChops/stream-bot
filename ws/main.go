package ws

import (
	"fmt"
	"time"

	"github.com/SnareChops/twitchbot/events"
	"nhooyr.io/websocket"
)

func Start(clientId, userId, token, wsUrl, subscribeUrl string, subscriptions []*events.EventSub) error {
	// Connect to ws server and get back session id
	var ws *websocket.Conn
	var sessionId string
	var err error
	for {
		ws, sessionId, err = connect(wsUrl)
		if err != nil {
			fmt.Printf("Failed to connect to twitch websocket: %s\n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	// Create subscriptions with session id
	// retry subscription creations that failed
	retry := subscriptions
	for len(retry) > 0 {
		retry = createSubscriptions(clientId, sessionId, userId, token, subscribeUrl, retry)
		if len(retry) > 0 {
			time.Sleep(5 * time.Second)
		}
	}

	// Listen for messages
	reconnect := make(chan string)
	close := make(chan bool)
	go listen(ws, reconnect, close, subscriptions)
	for {
		wsUrl = <-reconnect
		// Make ws connection
		prev := ws
		ws, _, err = connect(wsUrl)
		if err != nil {
			return err
		}
		// Cycle close channel
		prevClose := close
		close = make(chan bool)
		// Start new listener
		go listen(ws, reconnect, close, subscriptions)
		// Send close signal to previous listener
		prevClose <- true
		prev.Close(websocket.StatusGoingAway, "goodbye")
	}
}

func createSubscriptions(clientId, sessionId, userId, token, subscribeUrl string, subscriptions []*events.EventSub) []*events.EventSub {
	failed := []*events.EventSub{}
	for _, sub := range subscriptions {
		body := SubscribeBody{
			Type:      sub.Type,
			Version:   sub.Version,
			Condition: sub.Condition,
			Transport: SubscribeBodyTransport{
				Method:    "websocket",
				SessionID: sessionId,
			},
		}
		_, err := subscribe(clientId, token, subscribeUrl, body)
		if err != nil {
			fmt.Printf("Failed to create subscription for %s: %s\n", sub.Type, err)
			failed = append(failed, sub)
		} else {
			fmt.Printf("Subscribed to %s\n", sub.Type)
		}
	}
	return failed
}

type WelcomePayload struct {
	Session WelcomePayloadSession
}

type WelcomePayloadSession struct {
	ID               string
	Status           string
	ConnectedAt      string
	KeepaliveTimeout int
}
