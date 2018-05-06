package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/stjohnjohnson/reddit-watch/chatter"
	"github.com/stjohnjohnson/reddit-watch/data"
	"github.com/stjohnjohnson/reddit-watch/matcher"
	"github.com/stjohnjohnson/reddit-watch/scanner"
)

// These variables get set by the build script via the LDFLAGS
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type BotHandler struct {
	userData *data.UserData
	posts    scanner.ScannerChannel
	scan     *scanner.ScannerHandler
	messages chatter.ChatterChannel
	chat     *chatter.ChatterHandler
}

// /COMMAND OPTIONALDATA
var cmdRex = regexp.MustCompile(`(?i)^/(\w+)(?:\s(.+))?$`)

func (b *BotHandler) parseIncomingMessage(userID int64, message string) string {
	fields := cmdRex.FindStringSubmatch(message)
	if fields == nil {
		return ""
	}

	switch {
	case len(fields) > 2 && fields[1] == "buying":
		keyword := fields[2]

		err := b.userData.Add(userID, keyword)
		if err != nil {
			log.Println("Unable to add keyword: ", err)
		}

		return fmt.Sprintf("Okay, I'm going to watch for keywords matching %v", keyword)
	case len(fields) > 2 && fields[1] == "stop":
		keyword := fields[2]

		err := b.userData.Remove(userID, keyword)
		if err != nil {
			log.Println("Unable to remove keyword: ", err)
		}

		return fmt.Sprintf("I'm no longer going to watch for keywords matching %v", keyword)
	case len(fields) > 1 && fields[1] == "items":
		crit := b.userData.Get(userID)

		resp := []string{}
		for keyword, hits := range crit {
			resp = append(resp, fmt.Sprintf(" - %v (%d hits)", keyword, hits))
		}

		if len(resp) > 0 {
			return fmt.Sprintf("These are your current watch items:\n%v", strings.Join(resp, "\n"))
		} else {
			return "There are no items on your watch list"
		}

	case len(fields) > 1 && fields[1] == "stats":
		stats := b.userData.GetStats()

		resp := []string{
			"Interesting Statistics:",
		}

		for stat, val := range stats {
			resp = append(resp, fmt.Sprintf(" - %v (%v)", stat, val))
		}

		return strings.Join(resp, "\n")

	case len(fields) > 1 && fields[1] == "help":
		resp := []string{
			fmt.Sprintf("reddit-watcher@%v, commit %v, built at %v", version, commit, date),
			"https://git.io/vprdx",
			"How to use this bot:",
			" /buying <keyword> - will notify this channel if any listings match that keyword",
			" /stop <keyword> - will notify this channel if any listings match that keyword",
			" /items - returns list of watched items",
			" /stats - returns information about the bot",
			" /help - gets this help message",
		}

		return strings.Join(resp, "\n")

	default:
		return "That command doesn't look like anything to me."
	}
}

func (b *BotHandler) Loop() {
	for {
		select {
		case post := <-b.posts:
			forSale, mode, err := matcher.GetSale(post.Title)

			// Record stats
			b.userData.IncrementStat(mode)

			if err != nil {
				log.Printf("SKIP %s", err)
				continue
			}
			log.Printf("ACK %s", forSale)

			keywords := matcher.FindMatching(b.userData.GetKeywords(), forSale, post.SelfText)
			for _, keyword := range keywords {
				ids := b.userData.GetByKeyword(keyword)
				for _, id := range ids {
					log.Printf("MATCH %s (@%d)", keyword, id)
					b.chat.SendMessage(id, fmt.Sprintf("https://reddit.com%s (matching %s)", post.Permalink, keyword))
					b.userData.Increment(id, keyword)
				}
			}

		case update := <-b.messages:
			if update.Message == nil {
				continue
			}
			resp := b.parseIncomingMessage(update.Message.Chat.ID, update.Message.Text)

			if resp != "" {
				log.Printf("MSG [%s] %s", update.Message.From.UserName, update.Message.Text)
				b.chat.SendMessage(update.Message.Chat.ID, resp)
			}
		}
	}
}

func main() {
	log.Printf("Version %v, commit %v, built at %v", version, commit, date)

	token := flag.String("token", "INVALID", "Bot Token for Telegram")
	configPath := flag.String("config", "/config", "Location of user data")
	flag.Parse()

	userData, err := data.Load(*configPath)
	if err != nil {
		log.Printf("Unable to load config: %v", err)
	}

	scan, err := scanner.New(version)
	if err != nil {
		log.Fatalf("Failed to setup scanner: %v", err)
	}

	posts, err := scan.Start()
	if err != nil {
		log.Fatalf("Failed to start scanner: %v", err)
	}

	chat, err := chatter.New(version, *token)
	if err != nil {
		log.Fatalf("Failed to setup chatter: %v", err)
	}

	messages, err := chat.Start()
	if err != nil {
		log.Fatalf("Failed to start chatter: %v", err)
	}

	bot := &BotHandler{
		userData: userData,
		posts:    posts,
		scan:     scan,
		messages: messages,
		chat:     chat,
	}

	bot.Loop()
}
