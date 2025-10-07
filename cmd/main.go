package main

import (
	"discord-tasker-runner/config"
	"discord-tasker-runner/pkg/bot"
	"log"
	"os"
)

func main() {
	// parse config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// service account json check
	_, err = os.Stat("service-account.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("error checking service-account.json: file does not exist")
		} else {
			log.Fatal(err)
		}
	}

	// start bot
	bot.Run(*cfg)
}
