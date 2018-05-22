package bot

import (
	"fmt"
	"html"

	"github.com/stjohnjohnson/reddit-watcher/internal/matcher"
	"github.com/turnage/graw/reddit"
)

func (b *Handler) incomingPost(post *reddit.Post) error {
	item, err := matcher.ParseTitle(post.Title)
	if err != nil {
		return fmt.Errorf("unable to parse title: %s", err)
	}
	b.logger.Printf("PARSE: type: %s region: %s", item.Type, item.Region)

	// Record stats for type
	// @TODO Record stats for region
	err = b.stats.Increment(item.Type)
	if err != nil {
		return fmt.Errorf("unable to record stat: %s", err)
	}

	d, ok := b.data[item.Type]
	if !ok {
		return fmt.Errorf("unknown type: %s", item.Type)
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

			err = b.chat.SendMessage(id, fmt.Sprintf(`%s [<a href="%s">web</a>] [<a href="https://hidereferrer.com/?narwhal://open-url/www.reddit.com%s">app</a>] <i>(matched %s %s)</i>`,
				html.EscapeString(post.Title), post.URL, post.Permalink, item.Type, html.EscapeString(keyword)))
			if err != nil {
				b.logger.Printf("Unable to send message: %s", err)
			}

			err = d.Increment(id, keyword)
			if err != nil {
				b.logger.Printf("Unable to increment counter: %s", err)
			}
		}
	}

	return nil
}
