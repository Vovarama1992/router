package telegram

import (
	"context"

	"router/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
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
		b.app.API().Send(
			tgbotapi.NewMessage(chatID, "Ошибка создания конфига"),
		)
		return
	}

	doc := tgbotapi.NewDocument(chatID,
		tgbotapi.FileBytes{
			Name:  "wg.conf",
			Bytes: []byte(peer.Config),
		},
	)
	b.app.API().Send(doc)

	qr, _ := qrcode.Encode(peer.Config, qrcode.Medium, 256)
	photo := tgbotapi.NewPhoto(chatID,
		tgbotapi.FileBytes{Name: "wg.png", Bytes: qr},
	)
	b.app.API().Send(photo)
}
