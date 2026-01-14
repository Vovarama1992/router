package telegram

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	api *tgbotapi.BotAPI
}

func NewFromEnv() (*App, error) {
	log.Printf("[tg] NewFromEnv start")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Printf("[tg] TELEGRAM_BOT_TOKEN is empty")
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is empty")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("[tg] NewBotAPI FAILED err=%v", err)
		return nil, err
	}

	log.Printf("[tg] bot initialized username=%s", api.Self.UserName)
	return &App{api: api}, nil
}

func (a *App) Run(ctx context.Context, handler func(tgbotapi.Update)) {
	log.Printf("[tg] polling start")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := a.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[tg] polling stop (context done)")
			return
		case upd := <-updates:
			handler(upd)
		}
	}
}

func (a *App) API() *tgbotapi.BotAPI {
	return a.api
}
