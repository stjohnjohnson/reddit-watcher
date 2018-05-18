package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stjohnjohnson/reddit-watcher/matcher"
)

// /COMMAND OPTIONALDATA
var cmdRex = regexp.MustCompile(`(?i)^/(\w+)(?:\s(.+))?$`)

func (b *Handler) incomingMessage(userID int64, message string) {
	fields := cmdRex.FindStringSubmatch(message)
	if fields == nil {
		return
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
		b.logger.Printf("Unable to send message: %v", err)
	}
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

		return fmt.Sprintf("I'm no longer watching for %s posts that match %s", cmd, keyword)
	}

	err := b.data[cmd].Add(userID, keyword)
	if err != nil {
		b.logger.Println("Unable to add keyword: ", err)
	}

	// @TODO better message for ALL events
	return fmt.Sprintf("Okay, I'm going to watch for %s posts that match %s", cmd, keyword)
}

func (b *Handler) handleWatchlist(userID int64) string {
	resp := []string{}

	for _, t := range matcher.Types {
		crit := b.data[t].Get(userID)
		if len(crit) == 0 {
			continue
		}

		resp = append(resp, fmt.Sprintf("%s:", strings.ToUpper(t)))
		for keyword, hits := range crit {
			resp = append(resp, fmt.Sprintf(" - %v (%d hits)", keyword, hits))
		}
		resp = append(resp, "")
	}

	if len(resp) > 0 {
		return fmt.Sprintf("These are your current watch items:\n%v", strings.Join(resp, "\n"))
	}
	return "There are no items on your watch list"
}

func (b *Handler) handleHelp() string {
	resp := []string{
		fmt.Sprintf("Hi, I'm reddit-watcher@%v. I watch /r/mechmarket for specific keywords", b.version),
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
	}

	return strings.Join(resp, "\n")
}

func (b *Handler) handleStats() string {
	resp := []string{
		"Interesting Statistics:",
	}

	for stat, val := range b.stats.GetAll() {
		resp = append(resp, fmt.Sprintf(" - %v (%v)", stat, val))
	}

	return strings.Join(resp, "\n")
}
