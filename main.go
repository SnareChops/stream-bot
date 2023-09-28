package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/SnareChops/twitchbot/admin"
	"github.com/SnareChops/twitchbot/db"
	"github.com/SnareChops/twitchbot/events"
	"github.com/SnareChops/twitchbot/irc"
	"github.com/SnareChops/twitchbot/obs"
	"github.com/SnareChops/twitchbot/ui"
	"github.com/SnareChops/twitchbot/ws"
	"golang.org/x/exp/slices"
)

var TwitchScopes []string = []string{
	"chat:read",
	"chat:edit",
	events.CheerAuthScope,
	events.FollowAuthScope,
	events.RedeemAuthScope,
	events.SubscribeAuthScope,
}

func env(name string) string {
	result := os.Getenv(name)
	if result == "" {
		panic(name + " missing")
	}
	return result
}

func main() {
	clientId := env("TWITCH_CLIENT_ID")
	if slices.Contains(os.Args, "--token") {
		token(clientId)
		os.Exit(0)
	}
	channelName := env("TWITCH_CHANNEL_NAME")
	// clientSecret := env("TWITCH_CLIENT_SECRET")
	token := env("TWITCH_OAUTH_TOKEN")
	userId := env("TWITCH_USER_ID")
	wsUrl := env("TWITCH_WS_URL")
	apiUrl := env("TWITCH_API_URL")
	connString := env("MONGO_CONN_STRING")
	dbName := env("MONGO_DB_NAME")
	obsURL := env("OBS_URL")
	obsPassword := env("OBS_PASSWORD")

	uiOnly := flag.Bool("ui-only", false, "")
	flag.Parse()

	// Start the database connection to mongo
	disconnect := db.Start(connString, dbName)
	defer disconnect()
	if !*uiOnly {
		// Start the EventSub WS connection to twitch
		go ws.Start(clientId, userId, token, wsUrl, apiUrl, []*events.EventSub{
			events.NewCheerEventSub(userId),
			events.NewFollowEventSub(userId),
			events.NewGiftEventSub(userId),
			events.NewResubEventSub(userId),
			events.NewSubEventSub(userId),
			events.NewRedeemEventSub(userId),
		})

		// Start the IRC connection
		go irc.Start(channelName, token)
	}

	// Start the alert UI connection
	go ui.Start()

	// TODO: Start the admin interface connection
	go admin.Start()

	// Start OBS connection
	go obs.Start(obsURL, obsPassword)

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

func token(clientId string) {
	// cmd := exec.Command("twitch", "token", "--client-id", clientId, "-u", "-s", "chat:read chat:edit bits:read moderator:read:followers channel:read:redemptions channel:read:subscriptions")
	// stdout, err := cmd.Output()
	// if err != nil {
	// panic(err)
	// }
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(stdin) + "\n")
	regex := regexp.MustCompile(`(?m)User Access Token: (.*)`)
	match := regex.FindSubmatch(stdin)
	if len(match) == 0 || len(match[0]) == 0 {
		panic("No access token found in stdin")
	}
	filename := ".env"
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	content = regexp.MustCompile(`TWITCH_OAUTH_TOKEN=.*`).ReplaceAll(content, append([]byte("TWITCH_OAUTH_TOKEN="), match[1]...))
	err = os.WriteFile(filename, content, os.ModePerm)
	if err != nil {
		panic(err)
	}
	println("Token Updated")
}
