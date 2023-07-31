package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/SnareChops/twitchbot/events"
	"github.com/SnareChops/twitchbot/ui"
	"github.com/SnareChops/twitchbot/ws"
)

var TwitchScopes []string = []string{
	"chat:read",
	"chat:edit",
	events.CheerAuthScope,
	events.FollowAuthScope,
	events.RedeemAuthScope,
	events.SubscribeAuthScope,
}

func main() {
	// channel := "snarechops"
	clientId := os.Getenv("TWITCH_CLIENT_ID")
	if clientId == "" {
		panic("TWITCH_CLIENT_ID missing")
	}
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if clientSecret == "" {
		panic("TWITCH_CLIENT_SECRET missing")
	}
	token := os.Getenv("TWITCH_OAUTH_TOKEN")
	if token == "" {
		panic("TWITCH_OAUTH_TOKEN missing")
	}
	userId := os.Getenv("TWITCH_USER_ID")
	if userId == "" {
		panic("TWITCH_USER_ID")
	}
	wsUrl := os.Getenv("TWITCH_WS_URL")
	if wsUrl == "" {
		panic("TWITCH_WS_URL")
	}
	apiUrl := os.Getenv("TWITCH_API_URL")
	if apiUrl == "" {
		panic("TWITCH_API_URL")
	}

	send := make(chan []byte, 1)
	close := make(chan bool)

	// TODO: Send close signal when application is killed
	go ws.Start(send, clientId, userId, token, wsUrl, apiUrl, []*events.EventSub{
		events.NewCheerEventSub(userId),
		events.NewFollowEventSub(userId),
		events.NewGiftEventSub(userId),
		events.NewResubEventSub(userId),
		events.NewSubEventSub(userId),
	})

	// TODO: Send close signal when application is killed
	// go irc(channel, token)

	go ui.Start(send, close)

	// for {
	// 	time.Sleep(time.Second * 5)
	// 	send <- []byte("{\"kind\": \"twitch.follow\", \"username\": \"test\"}")
	// }

	println("Bot started...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	println("Interupt detected, closing bot...")
	// close <- true
	println("Goodbye...")
	os.Exit(0)
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getToken(id, secret string) (string, error) {
	data := url.Values{}
	data.Set("client_id", id)
	data.Set("client_secret", secret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("%s: %s", res.Status, string(content))
	}

	jsonResult := TokenResponse{}
	err = json.Unmarshal(content, &jsonResult)
	if err != nil {
		return "", err
	}

	return jsonResult.AccessToken, err
}
