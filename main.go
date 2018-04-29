package main

import (
	"log"

	"github.com/stjohnjohnson/reddit-watch/receiver"
	"github.com/stjohnjohnson/reddit-watch/scanner"
	"github.com/turnage/graw/reddit"
)

// These variables get set by the build script via the LDFLAGS
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log.Printf("Version %v, commit %v, built at %v", version, commit, date)

	posts := make(chan *reddit.Post)

	go scanner.Start(version, posts)
	receiver.Start(posts)

}
