package irc

import (
	"fmt"
	"strings"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

const (
	EVERYONE = 1 << iota
	VIP
	MODERATOR
)

func Start(channel, token string) error {
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
		handle(message)
	})
	go speaker(channel, client)
	return client.Connect()
}

func handle(message twitch.PrivateMessage) {
	content := strings.Split(strings.TrimSpace(message.Message), " ")
	if len(content) == 0 {
		return
	}
	if content[0][0] != '!' {
		scan(content)
		return
	}
	switch name := strings.ToLower(content[0]); name {
	default:
		command, err := FindCommand(name)
		if err != nil {
			fmt.Printf("Failed to find command: %s, %s", content[0], err)
		}
		if canRun(command.Permission, message) {
			signals.SendToIRC <- command.Response
		}
	}
}

func canRun(perm int, message twitch.PrivateMessage) bool {
	if _, ok := message.User.Badges["broadcaster"]; ok {
		return true
	}
	if perm&EVERYONE == EVERYONE {
		return true
	}
	if _, ok := message.User.Badges["vip"]; ok && perm&VIP == VIP {
		return true
	}
	if _, ok := message.User.Badges["moderator"]; ok && perm&MODERATOR == MODERATOR {
		return true
	}
	return false
}
