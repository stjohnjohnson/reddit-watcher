package bot

import (
	"fmt"

	"github.com/stjohnjohnson/reddit-watcher/matcher"
	"github.com/turnage/graw/reddit"
)

func (b *Handler) incomingPost(post *reddit.Post) {
	item, err := matcher.ParseTitle(post.Title)
	if err != nil {
		b.logger.Printf("Unable to parse title: %s", err)
		return
	}
	b.logger.Printf("PARSE: type: %s region: %s", item.Type, item.Region)

	// Record stats for type
	// @TODO Record stats for region
	err = b.stats.Increment(item.Type)
	if err != nil {
		b.logger.Printf("Unable to record stat: %s", err)
		return
	}

	d, ok := b.data[item.Type]
	if !ok {
		b.logger.Printf("Unknown type: %s", item.Type)
		return
	}

	keywords := matcher.FindMatching(d.GetKeywords(), item.Contents, post.SelfText)
	for _, keyword := range keywords {
		ids := d.GetByKeyword(keyword)
		for _, id := range ids {
			// @TODO Check region value
			if item.Region != "US" && item.Region != "" {
				continue
			}
			b.logger.Printf("MATCH: %s/%s for @%d, %s", item.Type, keyword, id, post.URL)

			err = b.chat.SendMessage(id, fmt.Sprintf("%s (matching %s %s)", post.URL, item.Type, keyword))
			if err != nil {
				b.logger.Printf("Unable to send message: %s", err)
			}

			err = d.Increment(id, keyword)
			if err != nil {
				b.logger.Printf("Unable to increment counter: %s", err)
			}
		}
	}
}
