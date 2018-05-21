package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/stjohnjohnson/reddit-watcher/chatter"
	"github.com/stjohnjohnson/reddit-watcher/data"
	"github.com/stjohnjohnson/reddit-watcher/matcher"
	"github.com/stjohnjohnson/reddit-watcher/scanner"
	"github.com/stjohnjohnson/reddit-watcher/stats"
)

// Handler is the bot object
type Handler struct {
	version  string
	data     map[string]data.Interface
	stats    stats.Interface
	posts    scanner.Channel
	scan     *scanner.Handler
	messages chatter.Channel
	chat     chatter.Interface
	logger   *log.Logger
}

// Loop is the main logic loop, listening for posts or messages from user
func (b *Handler) Loop() {
	for {
		select {
		case post := <-b.posts:
			b.logger.Printf("POST: %s", post.Title)
			err := b.incomingPost(post)
			if err != nil {
				b.logger.Printf("post failure: %v", err)
			}

		case update := <-b.messages:
			// Skip non-messages
			if update.Message == nil {
				continue
			}
			b.logger.Printf("MSG: %s: %s", update.Message.Chat.UserName, update.Message.Text)
			err := b.incomingMessage(update.Message.Chat.ID, update.Message.Text)
			if err != nil {
				b.logger.Printf("message failure: %v", err)
			}
		}
	}
}

// New creates a new bot given a Telegram token and config directory
func New(token, configDir, version string) (*Handler, error) {
	appData := make(map[string]data.Interface)
	logger := log.New(os.Stderr, "[BOT]  ", log.LstdFlags)

	for _, t := range matcher.Types {
		d, err := data.Load(fmt.Sprintf("%s/%s", configDir, t))
		if err != nil {
			logger.Printf("Unable to load config: %v", err)
		}
		appData[t] = d
	}

	statData, err := stats.Load(configDir)
	if err != nil {
		logger.Printf("Unable to load stats: %v", err)
	}

	scan, err := scanner.New(version)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup scanner: %v", err)
	}

	posts, err := scan.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start scanner: %v", err)
	}

	chat, err := chatter.New(version, token)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup chatter: %v", err)
	}

	messages, err := chat.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start chatter: %v", err)
	}

	return &Handler{
		version:  version,
		data:     appData,
		stats:    statData,
		posts:    posts,
		scan:     scan,
		messages: messages,
		chat:     chat,
		logger:   logger,
	}, nil
}
