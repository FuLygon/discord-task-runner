package main

import (
	"discord-tasker-runner/config"
	"discord-tasker-runner/pkg/bot"
	"log"
)

func main() {
	// parse config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// start bot
	bot.Run(*cfg)
}
