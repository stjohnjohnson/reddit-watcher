package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stjohnjohnson/reddit-watch/matcher"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

// These variables get set by the build script via the LDFLAGS
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type monitorScript struct {
	script reddit.Script
}

func (r *monitorScript) Post(p *reddit.Post) error {
	forSale, err := matcher.FindValidSale(p.Title)

	if err == nil {
		fmt.Printf("MATCHED %s | https://reddit.com%s\n", forSale, p.Permalink)
	} else {
		fmt.Printf("SKIPPED %s\n", err)
	}

	return nil
}

func main() {
	log.Printf("Version %v, commit %v, built at %v", version, commit, date)

	rate := time.Second * 5
	script, _ := reddit.NewScript(fmt.Sprintf("golang:reddit-watch:%v (by /u/GalacticGargleBlaster)", version), rate)

	cfg := graw.Config{Subreddits: []string{"mechmarket"}}
	handler := &monitorScript{script: script}

	if _, wait, err := graw.Scan(handler, script, cfg); err != nil {
		log.Println("Failed to start graw run: ", err)
	} else {
		log.Println("graw run failed: ", wait())
	}
}
