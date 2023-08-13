package irc

import (
	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

func speaker(channel string, client *twitch.Client) {
	for {
		select {
		case <-signals.Shutdown:
			return
		case message := <-signals.SendToIRC:
			client.Say(channel, message)
		}
	}
}
