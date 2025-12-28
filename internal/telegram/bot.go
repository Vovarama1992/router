package telegram

import (
	"context"
	"log"

	"router/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	app *App
	svc *domain.Service
}

func NewBot(app *App, svc *domain.Service) *Bot {
	return &Bot{
		app: app,
		svc: svc,
	}
}

func (b *Bot) Handle(update tgbotapi.Update) {
	if update.Message != nil {
		b.onMessage(update.Message)
		return
	}

}

func (b *Bot) onMessage(msg *tgbotapi.Message) {
	if msg.Text == "/start" {
		m := tgbotapi.NewMessage(
			msg.Chat.ID,
			"Нажми кнопку ниже, чтобы получить VPN-конфиг",
		)
		m.ReplyMarkup = mainKeyboard()
		b.app.API().Send(m)
		return
	}

	if msg.Text == "Получить конфиг" {
		b.sendConfig(msg.Chat.ID)
	}
}

func (b *Bot) sendConfig(chatID int64) {
	peer, err := b.svc.CreatePeer(context.Background())
	if err != nil {
		log.Println("[BOT] CreatePeer error:", err)
		b.app.API().Send(
			tgbotapi.NewMessage(chatID, "Ошибка создания конфига"),
		)
		return
	}

	// 1. файл
	doc := tgbotapi.NewDocument(chatID,
		tgbotapi.FileBytes{
			Name:  "client.ovpn",
			Bytes: []byte(peer.Config),
		},
	)
	if _, err := b.app.API().Send(doc); err != nil {
		log.Println("[BOT] send file error:", err)
		return
	}
}
