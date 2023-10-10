package irc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SnareChops/twitchbot/signals"
	"github.com/gempir/go-twitch-irc/v4"
)

type ResubSignal struct {
	Kind     string  `json:"kind"`
	Username string  `json:"username"`
	Tier     int     `json:"tier"`
	Months   int     `json:"months"`
	Streak   int     `json:"streak"`
	Message  string  `json:"message"`
	Emotes   []Emote `json:"emotes"`
}

type Emote struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Count     int             `json:"count"`
	Positions []EmotePosition `json:"position"`
}

type EmotePosition struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func resub(message twitch.UserNoticeMessage) {
	months, err := strconv.ParseInt(message.MsgParams["msg-param-cumulative-months"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse cumulative months: %s\n", err.Error())
	}
	streak, err := strconv.ParseInt(message.MsgParams["msg-param-streak-months"], 10, 0)
	if err != nil {
		fmt.Printf("Failed to parse streak: %s\n", err.Error())
	}

	signal, err := json.Marshal(ResubSignal{
		Kind:     "twitch.resub",
		Username: message.User.DisplayName,
		Tier:     transformTier(message.MsgParams["msg-param-sub-plan"]),
		Months:   int(months),
		Streak:   int(streak),
		Message:  message.Message,
		Emotes:   transformEmotes(message.Emotes),
	})
	if err != nil {
		fmt.Printf("Failed to marshal ResubSignal: %s\n", err.Error())
		return
	}
	signals.SendToUI <- signal
}

func transformEmotes(emotes []*twitch.Emote) []Emote {
	result := []Emote{}
	for _, emote := range emotes {
		result = append(result, Emote{
			ID:        emote.ID,
			Name:      emote.Name,
			Count:     emote.Count,
			Positions: transformPositions(emote.Positions),
		})
	}
	return result
}

func transformPositions(positions []twitch.EmotePosition) []EmotePosition {
	result := []EmotePosition{}
	for _, position := range positions {
		result = append(result, EmotePosition{
			Start: position.Start,
			End:   position.End,
		})
	}
	return result
}
