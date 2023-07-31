package main

import (
	"github.com/gempir/go-twitch-irc/v4"
)

func irc(channel, token string) error {
	println("Connecting to IRC")
	client := twitch.NewClient(channel, "oauth:"+token)

	client.Capabilities = []string{twitch.TagsCapability, twitch.CommandsCapability}
	client.Join(channel)
	client.OnConnect(func() {
		println("Connected to IRC")
	})
	client.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		println("Reconnected to IRC")
	})
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {

	})
	client.Say(channel, "test message")
	return client.Connect()
}
