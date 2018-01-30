package main

import (
	"os"

	"github.com/botopolis/bot"
	"github.com/botopolis/slack"
	"github.com/botopolis/redis"
)

func main() {
	bot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
		redis.New(os.Getenv("REDIS_URL")),
		welcomer{},
		webhook{},
	).Run()
}
