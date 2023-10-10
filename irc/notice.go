package irc

import "github.com/gempir/go-twitch-irc/v4"

func notice(message twitch.UserNoticeMessage) {
	switch message.MsgID {
	case "sub":
		sub(message)
	case "resub":
		resub(message)
	case "subgift":
		subgift(message)
	case "submysterygift":
		submysterygift(message)
	case "primepaidupgrade":
		primepaidupgrade(message)
	case "giftpaidupgrade":
		giftpaidupgrade(message)
	case "anongiftpaidupgrade":
		anongiftpaidupgrade(message)
	case "raid":
		raid(message)
	case "ritual":
		ritual(message)
	case "bitsbadgetier":
		bitsbadgetier(message)
	}
}
