package obs

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SnareChops/twitchbot/signals"
	"nhooyr.io/websocket"
)

const (
	OpCodeHello                = 0
	OpCodeIdentify             = 1
	OpCodeIdentified           = 2
	OpCodeReidentify           = 3
	OpCodeEvent                = 5
	OpCodeRequest              = 6
	OpCodeRequestResponse      = 7
	OpCodeRequestBatch         = 8
	OpCodeRequestBatchResponse = 9
)

type OBSMessage[T any] struct {
	OP   int `json:"op"`
	Data T   `json:"d"`
}

func NewOBSMessage[T any](op int, data T) *OBSMessage[T] {
	return &OBSMessage[T]{
		OP:   op,
		Data: data,
	}
}

type Client struct {
	id            int
	conn          *websocket.Conn
	password      string
	version       string
	challenge     string
	salt          string
	token         string
	subscriptions map[int]func(message []byte)
}

func (self *Client) Connect(url, password string) error {
	self.password = password
	fmt.Printf("Connecting OBS Websocket to %s\n", url)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var err error
	self.conn, _, err = websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}
	err = self.waitForHello()
	if err != nil {
		return err
	}
	err = self.identify()
	if err != nil {
		return err
	}
	err = self.waitForIdentified()
	if err != nil {
		return err
	}
	return err
}

func (self *Client) Subscribe(id int, handler func([]byte)) {
	if self.subscriptions == nil {
		self.subscriptions = map[int]func([]byte){}
	}
	self.subscriptions[id] = handler
}

type EventMessage[T any] struct {
	Type   string `json:"eventType"`
	Intent int    `json:"eventIntent"`
	Data   T      `json:"eventData"`
}

func (self *Client) Listen() {
	for {
		select {
		case <-signals.Shutdown:
			return
		default:
			_, data, err := self.conn.Read(context.Background())
			if err != nil {
				fmt.Printf("Failed to read OBS websocket message: %s\n", err.Error())
			}
			message := &OBSMessage[any]{}
			err = json.Unmarshal(data, message)
			if err != nil {
				fmt.Printf("Failed to decode OBS websocket message: %s\n", err.Error())
			}
			switch message.OP {
			case OpCodeEvent:
				event := &OBSMessage[EventMessage[any]]{}
				err = json.Unmarshal(data, event)
				if err != nil {
					fmt.Printf("Failed to decode OBS websocket event: %s\n", err.Error())
				}
				if handler, ok := self.subscriptions[event.Data.Intent]; ok {
					go handler(data)
				}
			default:
			}
		}
	}
}

type HelloMessage struct {
	SocketVersion  string `json:"obsWebSocketVersion"`
	RPCVersion     int    `json:"rpcVersion"`
	Authentication struct {
		Challenge string `json:"challenge"`
		Salt      string `json:"salt"`
	} `json:"authentication"`
}

func (self *Client) waitForHello() error {
	_, data, err := self.conn.Read(context.Background())
	if err != nil {
		return err
	}
	hello := &OBSMessage[HelloMessage]{}
	err = json.Unmarshal(data, hello)
	if err != nil {
		return err
	}
	self.version = hello.Data.SocketVersion
	self.challenge = hello.Data.Authentication.Challenge
	self.salt = hello.Data.Authentication.Salt
	err = self.encodeAuth()
	if err != nil {
		return err
	}
	return nil
}

func (self *Client) encodeAuth() error {
	if self.challenge == "" && self.salt == "" {
		return nil
	}
	// Create SHA256 hash of the password + salt
	secret, err := sha256AndBase64Encode([]byte(self.password + self.salt))
	if err != nil {
		return err
	}
	self.token, err = sha256AndBase64Encode([]byte(secret + self.challenge))
	if err != nil {
		return err
	}
	return nil
}

type IdentifyMessage struct {
	RPCVersion     int    `json:"rpcVersion"`
	Authentication string `json:"authentication"`
	Subscriptions  int    `json:"eventSubscriptions"`
}

func (self *Client) identify() error {
	subscriptions := 0
	for key, _ := range self.subscriptions {
		subscriptions |= key
	}
	message := NewOBSMessage(OpCodeIdentify, IdentifyMessage{
		RPCVersion:     1,
		Authentication: self.token,
		Subscriptions:  subscriptions,
	})
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return self.conn.Write(context.Background(), websocket.MessageText, data)
}

func (self *Client) waitForIdentified() error {
	_, data, err := self.conn.Read(context.Background())
	if err != nil {
		return err
	}
	message := &OBSMessage[any]{}
	err = json.Unmarshal(data, message)
	if err != nil {
		return err
	}
	if message.OP != OpCodeIdentified {
		return fmt.Errorf("Expected identified but received other message: %v\n", message)
	}
	return nil
}

func sha256AndBase64Encode(data []byte) (string, error) {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}
