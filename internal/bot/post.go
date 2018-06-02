package bot

import (
	"fmt"
	"html"
	"regexp"

	"github.com/stjohnjohnson/reddit-watcher/internal/matcher"
	"github.com/turnage/graw/reddit"
)

var messageTemplate = `%s [<a href="%s">web</a>] [<a href="https://git.io/vhZZN#%s">app</a>] <i>(matched %s %s)</i>`

func (b *Handler) incomingPost(post *reddit.Post) error {
	item, err := matcher.ParseTitle(post.Title)
	if err != nil {
		return fmt.Errorf("unable to parse title: %s", err)
	}
	b.logger.Printf("PARSE: type: %s region: %s", item.Type, item.Region)

	// Record stats for type
	// @TODO Record stats for region
	b.stats.Increment(item.Type)

	d, ok := b.data[item.Type]
	if !ok {
		return fmt.Errorf("unknown type: %s", item.Type)
	}

	keywords := matcher.FindMatching(d.GetKeywords(), item.Contents, post.SelfText)
	for _, keyword := range keywords {
		escapedKeyword := html.EscapeString(keyword)
		escapedTitle := html.EscapeString(post.Title)
		keywordReplacer := regexp.MustCompile(`(?i)(\[[^\]]+\])`)
		if keyword != "*" {
			keywordReplacer = regexp.MustCompile("(?i)(" + regexp.QuoteMeta(escapedKeyword) + ")")
		}
		escapedTitle = keywordReplacer.ReplaceAllString(escapedTitle, "<b>$1</b>")
		message := fmt.Sprintf(messageTemplate, escapedTitle, post.URL, post.Permalink, item.Type, escapedKeyword)

		ids := d.GetByKeyword(keyword)
		for _, id := range ids {
			// @TODO Check region value
			if item.Region != "US" && item.Region != "" {
				continue
			}
			b.logger.Printf("MATCH: %s/%s for @%d, %s", item.Type, keyword, id, post.URL)

			err = b.chat.SendMessage(id, message)
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
