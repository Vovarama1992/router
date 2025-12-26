package telegram

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	api *tgbotapi.BotAPI
}

func NewFromEnv() (*App, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is empty")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &App{api: api}, nil
}

func (a *App) Run(ctx context.Context, handler func(tgbotapi.Update)) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := a.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case upd := <-updates:
			handler(upd)
		}
	}
}

func (a *App) API() *tgbotapi.BotAPI {
	return a.api
}
