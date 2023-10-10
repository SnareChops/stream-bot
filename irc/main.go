package irc

import (
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
	client.OnConnect(connect)
	client.OnReconnectMessage(reconnect)
	client.OnPrivateMessage(private)
	client.OnUserJoinMessage(join)
	client.OnUserPartMessage(leave)
	// client.OnUserNoticeMessage(notice)
	go speaker(channel, client)
	return client.Connect()
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
