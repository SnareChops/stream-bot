package obs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SnareChops/twitchbot/signals"
)

func Start(url, password string) {
	client := new(Client)
	client.Subscribe(1<<16, InputVolumeMetersHandler)
	var err error
	for {
		err = client.Connect(url, password)
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to obs websocket: %s\n", err.Error())
		time.Sleep(5 * time.Second)
	}
	client.Listen()
}

type InputVolumeMetersMessage struct {
	Inputs []struct {
		Name   string      `json:"inputName"`
		Levels [][]float64 `json:"inputLevelsMul"`
	} `json:"inputs"`
}

type AvatarSignal struct {
	Kind    string `json:"kind"`
	Type    string `json:"type"`
	Talking bool   `json:"talking"`
}

func InputVolumeMetersHandler(message []byte) {
	event := &OBSMessage[EventMessage[InputVolumeMetersMessage]]{}
	err := json.Unmarshal(message, event)
	if err != nil {
		fmt.Printf("Failed to decode OBS InputVolumeMeters event: %s\n", err.Error())
	}
	for _, input := range event.Data.Data.Inputs {
		if input.Name == "Mic" {
			data, err := json.Marshal(AvatarSignal{
				Kind:    "avatar",
				Type:    "talking",
				Talking: input.Levels[0][0] > .001,
			})
			if err != nil {
				fmt.Printf("Failed to encode AvatarSignal: %s\n", err.Error())
			}
			signals.SendToUI <- data
			break
		}
	}
}
