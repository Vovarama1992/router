package telegram

import (
	"context"

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

	if update.CallbackQuery != nil {
		b.onCallback(update.CallbackQuery)
		return
	}
}

func (b *Bot) onMessage(msg *tgbotapi.Message) {
	if msg.Text != "/start" {
		return
	}

	m := tgbotapi.NewMessage(
		msg.Chat.ID,
		"Нажми кнопку, чтобы получить VPN-конфиг",
	)
	m.ReplyMarkup = mainKeyboard()

	b.app.API().Send(m)
}

func (b *Bot) onCallback(cb *tgbotapi.CallbackQuery) {
	b.app.API().Request(tgbotapi.NewCallback(cb.ID, ""))

	if cb.Data != "get_config" {
		return
	}

	peer, err := b.svc.CreatePeer(context.Background())
	if err != nil {
		b.app.API().Send(
			tgbotapi.NewMessage(cb.Message.Chat.ID, "Ошибка создания конфига"),
		)
		return
	}

	doc := tgbotapi.NewDocument(
		cb.Message.Chat.ID,
		tgbotapi.FileBytes{
			Name:  "wg.conf",
			Bytes: []byte(peer.Config),
		},
	)

	b.app.API().Send(doc)
}
