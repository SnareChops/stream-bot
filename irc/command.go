package irc

import (
	"context"
	"time"

	"github.com/SnareChops/twitchbot/db"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CommandCollection = "command"
)

type Command struct {
	Name       string `bson:"name"`
	Permission int    `bson:"permission"`
	Response   string `bson:"response"`
}

func FindCommand(name string) (*Command, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	command := new(Command)
	err := db.Database.Collection(CommandCollection).FindOne(ctx, bson.D{{"name", name}}).Decode(command)
	return command, err
}
