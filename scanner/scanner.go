package scanner

import (
	"errors"
	"fmt"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type ScannerHandler struct {
	script  reddit.Script
	config  graw.Config
	channel ScannerChannel
}

type ScannerChannel chan *reddit.Post

func (r *ScannerHandler) Post(p *reddit.Post) error {
	r.channel <- p

	return nil
}

// Start will start the scanner and return a channel to listen for new posts
func (r *ScannerHandler) Start() (chan *reddit.Post, error) {
	_, _, err := graw.Scan(r, r.script, r.config)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to start: %v", err))
	}

	return r.channel, err
}

// New creates a new scanner to look for Reddit posts in /r/mechmarket
func New(version string) (*ScannerHandler, error) {
	script, err := reddit.NewScript(
		fmt.Sprintf("golang:reddit-watch:%v (by /u/GalacticGargleBlaster)", version),
		time.Second*5,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to setup: %v", err))
	}

	channel := make(ScannerChannel)
	config := graw.Config{
		Subreddits: []string{"mechmarket"},
	}

	return &ScannerHandler{
		script:  script,
		config:  config,
		channel: channel,
	}, nil
}
