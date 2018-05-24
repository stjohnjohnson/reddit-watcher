package chatter

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

// Handler is a telegram bot
type Handler struct {
	bot    *tgbotapi.BotAPI
	logger *log.Logger
}

// Interface is the stats public functions
type Interface interface {
	Start() (Channel, error)
	SendMessage(int64, string) error
}

// Channel is a message channel
type Channel tgbotapi.UpdatesChannel

// Start begins listening to messages from Telegram
func (r *Handler) Start() (Channel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	c, err := r.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, fmt.Errorf("Unable to start: %v", err)
	}

	return Channel(c), nil
}

// SendMessage will send a message to a given user
func (r *Handler) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	_, err := r.bot.Send(msg)

	if err != nil {
		return fmt.Errorf("Unable to send: %v", err)
	}
	return nil
}

// New creates a new Telegram bot
func New(version, token string) (*Handler, error) {
	logger := log.New(os.Stderr, "[CHAT] ", log.LstdFlags)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup: %v", err)
	}

	logger.Printf("Authorized on account %s", bot.Self.UserName)

	return &Handler{
		bot:    bot,
		logger: logger,
	}, nil
}
