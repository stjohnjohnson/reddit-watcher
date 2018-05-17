package scanner

import (
	"fmt"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

// Handler is a reddit bot
type Handler struct {
	script  reddit.Script
	config  graw.Config
	channel Channel
}

// Channel is a reddit post channel
type Channel chan *reddit.Post

// Post receives an update from Reddit and forwards it to the Channel
func (r *Handler) Post(p *reddit.Post) error {
	r.channel <- p

	return nil
}

// Start will start the scanner and return a channel to listen for new posts
func (r *Handler) Start() (chan *reddit.Post, error) {
	_, _, err := graw.Scan(r, r.script, r.config)

	if err != nil {
		return nil, fmt.Errorf("Unable to start: %v", err)
	}

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

	channel := make(Channel)
	config := graw.Config{
		Subreddits: []string{"mechmarket"},
	}

	return &Handler{
		script:  script,
		config:  config,
		channel: channel,
	}, nil
}
