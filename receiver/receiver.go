package receiver

import (
	"log"

	"github.com/stjohnjohnson/reddit-watch/matcher"
	"github.com/turnage/graw/reddit"
)

var criteria = [...]string{
	"demonic",
	"sanctuary",
	"tai hao",
	"tai-hao",
	"taihao",
}

func Start(channel chan *reddit.Post) {
	for {
		p, ok := <-channel

		if !ok {
			return
		}

		forSale, err := matcher.FindValidSale(p.Title)

		if err == nil {
			log.Printf("MATCHED %s | https://reddit.com%s\n", forSale, p.Permalink)
		} else {
			log.Printf("SKIPPED %s\n", err)
		}
	}
}
