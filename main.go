package main

import (
	"flag"
	"log"

	"github.com/stjohnjohnson/reddit-watcher/bot"
)

// These variables get set by the build script via the LDFLAGS
var (
	version = "0.0.0-dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log.SetPrefix("[MAIN] ")
	log.Printf("Version %v, commit %v, built at %v", version, commit, date)

	token := flag.String("token", "INVALID", "Bot Token for Telegram")
	configPath := flag.String("config", "/config", "Location of user data")
	flag.Parse()

	bot, err := bot.New(*token, *configPath, version)
	if err != nil {
		log.Fatalf("Unable to start bot: %v", err)
	}

	bot.Loop()
}
