package telegram

import (
	"context"
	"log"

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
	log.Println("[BOT] sendConfig called")

	peer, err := b.svc.CreatePeer(context.Background())
	if err != nil {
		log.Printf("[BOT] CreatePeer error: %v", err)

		b.app.API().Send(
			tgbotapi.NewMessage(chatID, "Ошибка создания конфига"),
		)
		return
	}

	log.Println("[BOT] Config generated, sending file")

	doc := tgbotapi.NewDocument(chatID,
		tgbotapi.FileBytes{
			Name:  "client.ovpn",
			Bytes: []byte(peer.Config),
		},
	)
	b.app.API().Send(doc)

	qr, err := qrcode.Encode(peer.Config, qrcode.Medium, 256)
	if err != nil {
		log.Printf("[BOT] QR error: %v", err)
		return
	}

	photo := tgbotapi.NewPhoto(chatID,
		tgbotapi.FileBytes{Name: "client.png", Bytes: qr},
	)
	b.app.API().Send(photo)

	log.Println("[BOT] Config sent successfully")
}
