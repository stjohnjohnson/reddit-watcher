package chatter

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

type ChatterHandler struct {
	bot *tgbotapi.BotAPI
}

type ChatterChannel tgbotapi.UpdatesChannel

func (r *ChatterHandler) Start() (ChatterChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	c, err := r.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to start: %v", err))
	}

	return ChatterChannel(c), nil
}

func (r *ChatterHandler) SendMessage(chatID int64, message string) error {
	_, err := r.bot.Send(tgbotapi.NewMessage(chatID, message))

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to send: %v", err))
	}
	return nil
}

func New(version, token string) (*ChatterHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to setup: %v", err))
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &ChatterHandler{
		bot: bot,
	}, nil
}
