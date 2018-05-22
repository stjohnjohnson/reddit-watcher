package scanner

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

// Handler is a reddit bot
type Handler struct {
	script  reddit.Script
	config  graw.Config
	channel Channel
	logger  *log.Logger
}

// Interface is the stats public functions
type Interface interface {
	Post(*reddit.Post) error
	Start() (chan *reddit.Post, error)
}

// Channel is a reddit post channel
type Channel chan *reddit.Post

// Post receives an update from Reddit and forwards it to the Channel
func (r *Handler) Post(p *reddit.Post) error {
	r.logger.Printf("Received post: %v", p.URL)
	r.channel <- p

	return nil
}

// Start will start the scanner and return a channel to listen for new posts
func (r *Handler) Start() (chan *reddit.Post, error) {
	_, wait, err := graw.Scan(r, r.script, r.config)

	if err != nil {
		r.logger.Printf("Scanner failed to start: %v", err)
		return nil, fmt.Errorf("Unable to start: %v", err)
	}
	r.logger.Print("Scanner started")

	// @TODO Add retries instead of exiting the app
	go func() {
		r.logger.Fatalf("Scanner failed: %v", wait())
	}()

	return r.channel, err
}

// New creates a new scanner to look for Reddit posts in /r/mechmarket
func New(version string) (*Handler, error) {
	script, err := reddit.NewScript(
		fmt.Sprintf("golang:reddit-watcher:%v (by /u/GalacticGargleBlaster)", version),
		time.Second*15,
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup: %v", err)
	}

	logger := log.New(os.Stderr, "[SCAN] ", log.LstdFlags)
	channel := make(Channel)
	config := graw.Config{
		Subreddits: []string{"mechmarket"},
		Logger:     logger,
	}

	return &Handler{
		script:  script,
		config:  config,
		channel: channel,
		logger:  logger,
	}, nil
}
