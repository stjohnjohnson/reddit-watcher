package bot

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"

	"github.com/stjohnjohnson/reddit-watcher/internal/matcher"
)

var helpText = strings.Join([]string{
	"",
	"",
	"You can subscribe or unsubscribe to events by using the following commands:",
	" /buying <keyword> - something being bought",
	" /selling <keyword> - something being sold",
	" /vendor <keyword> - updates from vendors",
	" /artisan <keyword> - updates from artisans",
	" /groupbuy <keyword> - updates about group buys",
	" /interestcheck <keyword> - feedback about a design",
	" /giveaway <keyword> - something being given away",
	"",
	"Other options:",
	" /items - returns list of watched items",
	" /help - gets this help message",
}, "\n")

// /COMMAND OPTIONALDATA
var cmdRex = regexp.MustCompile(`(?i)^/(\w+)(?:\s(.+))?$`)

func (b *Handler) incomingMessage(userID int64, message string) error {
	fields := cmdRex.FindStringSubmatch(message)
	if fields == nil {
		return nil
	}

	var resp string
	switch cmd := fields[1]; cmd {
	case matcher.Buying, matcher.Selling, matcher.Artisan, matcher.Vendor,
		matcher.GroupBuy, matcher.InterestCheck, matcher.Giveaway:
		resp = b.handleSubscribe(userID, cmd, fields[2])

	case "items":
		resp = b.handleWatchlist(userID)

	case "stats":
		resp = b.handleStats()

	case "help":
		resp = b.handleHelp()

	default:
		resp = "That command doesn't look like anything to me."
	}

	err := b.chat.SendMessage(userID, resp)
	if err != nil {
		return fmt.Errorf("Unable to send message: %v", err)
	}

	return nil
}

func (b *Handler) handleSubscribe(userID int64, cmd, keyword string) string {
	if keyword == "" {
		keyword = "*"
	}

	if b.data[cmd].Exists(userID, keyword) {
		err := b.data[cmd].Remove(userID, keyword)
		if err != nil {
			b.logger.Println("Unable to remove keyword: ", err)
		}

		return fmt.Sprintf("I'm no longer watching for <b>%s</b> posts that match <b>%s</b>", html.EscapeString(cmd), html.EscapeString(keyword))
	}

	err := b.data[cmd].Add(userID, keyword)
	if err != nil {
		b.logger.Println("Unable to add keyword: ", err)
	}

	// @TODO better message for ALL events
	return fmt.Sprintf("Okay, I'm going to watch for <b>%s</b> posts that match <b>%s</b>", html.EscapeString(cmd), html.EscapeString(keyword))
}

func (b *Handler) handleWatchlist(userID int64) string {
	resp := []string{}

	for _, t := range matcher.Types {
		crit := b.data[t].Get(userID)
		if len(crit) == 0 {
			continue
		}

		keys := make([]string, 0)
		for k := range crit {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		resp = append(resp, fmt.Sprintf("%s:", strings.ToUpper(t)))
		for _, keyword := range keys {
			resp = append(resp, fmt.Sprintf(" - %v <i>(%d hits)</i>", html.EscapeString(keyword), crit[keyword]))
		}
		resp = append(resp, "")
	}

	if len(resp) > 0 {
		return fmt.Sprintf("These are your current watch items:\n%v", strings.Join(resp, "\n"))
	}
	return "There are no items on your watch list"
}

func (b *Handler) handleHelp() string {
	return fmt.Sprintf(`Hi, I'm <a href="https://github.com/stjohnjohnson/reddit-watcher">reddit-watcher@%v</a>. I watch /r/mechmarket for specific keywords%s`, b.version, html.EscapeString(helpText))
}

func (b *Handler) handleStats() string {
	resp := []string{
		"Interesting Statistics:",
	}

	for stat, val := range b.stats.GetAll() {
		resp = append(resp, fmt.Sprintf(" - %v <i>(%v)</i>", stat, val))
	}

	return strings.Join(resp, "\n")
}
