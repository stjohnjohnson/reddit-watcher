package scanner

import (
	"fmt"
	"log"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type scannerHandler struct {
	script  reddit.Script
	channel chan *reddit.Post
}

func (r *scannerHandler) Post(p *reddit.Post) error {
	r.channel <- p

	return nil
}

func Start(version string, channel chan *reddit.Post) {
	script, _ := reddit.NewScript(
		fmt.Sprintf("golang:reddit-watch:%v (by /u/GalacticGargleBlaster)", version),
		time.Second*5,
	)
	cfg := graw.Config{
		Subreddits: []string{"mechmarket"},
	}
	handler := &scannerHandler{
		script:  script,
		channel: channel,
	}

	if _, wait, err := graw.Scan(handler, script, cfg); err != nil {
		log.Println("Failed to start: ", err)
	} else {
		log.Println("Failed to scan: ", wait())
	}
}
