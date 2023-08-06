package ws

import (
	"fmt"
	"time"

	"github.com/SnareChops/twitchbot/events"
	"nhooyr.io/websocket"
)

func Start(send chan []byte, clientId, userId, token, wsUrl, subscribeUrl string, subscriptions []*events.EventSub) error {
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
	failed := subscriptions
	for len(failed) > 0 {
		failed = createSubscriptions(clientId, sessionId, userId, token, subscribeUrl, subscriptions)
		if len(failed) > 0 {
			time.Sleep(5 * time.Second)
		}
	}

	// Listen for messages
	reconnect := make(chan string)
	close := make(chan bool)
	go listen(ws, send, reconnect, close, subscriptions)
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
		go listen(ws, send, reconnect, close, subscriptions)
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
			fmt.Printf("Failed to create subscription: %s\n", err)
			failed = append(failed, sub)
		}
		fmt.Printf("Subscribed to %s\n", sub.Type)
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
