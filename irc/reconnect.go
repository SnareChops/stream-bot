package irc

import "github.com/gempir/go-twitch-irc/v4"

func reconnect(message twitch.ReconnectMessage) {
	println("Reconnected to IRC")
}
