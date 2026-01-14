package telegram

import (
	"context"
	"log"
	"time"

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
	log.Printf(
		"[tg] update received hasMessage=%v",
		update.Message != nil,
	)

	if update.Message != nil {
		b.onMessage(update.Message)
		return
	}
}

func (b *Bot) onMessage(msg *tgbotapi.Message) {
	log.Printf(
		"[tg] message chatID=%d user=%s text=%q",
		msg.Chat.ID,
		msg.From.UserName,
		msg.Text,
	)

	if msg.Text == "/start" {
		log.Printf("[tg] /start chatID=%d", msg.Chat.ID)

		m := tgbotapi.NewMessage(
			msg.Chat.ID,
			"Нажми кнопку ниже, чтобы получить VPN-конфиг",
		)
		m.ReplyMarkup = mainKeyboard()

		if _, err := b.app.API().Send(m); err != nil {
			log.Printf("[tg] send /start FAILED chatID=%d err=%v", msg.Chat.ID, err)
		}
		return
	}

	if msg.Text == "Получить конфиг" {
		log.Printf("[tg] get config pressed chatID=%d", msg.Chat.ID)
		b.sendConfig(msg.Chat.ID)
	}
}

func (b *Bot) sendConfig(chatID int64) {
	start := time.Now()

	log.Printf("[tg] sendConfig start chatID=%d", chatID)

	peer, err := b.svc.CreatePeer(context.Background())
	if err != nil {
		log.Printf("[tg] CreatePeer FAILED chatID=%d err=%v", chatID, err)
		b.app.API().Send(
			tgbotapi.NewMessage(chatID, "Ошибка создания конфига"),
		)
		return
	}

	log.Printf("[tg] CreatePeer OK chatID=%d configBytes=%d", chatID, len(peer.Config))

	doc := tgbotapi.NewDocument(
		chatID,
		tgbotapi.FileBytes{
			Name:  "client.ovpn",
			Bytes: []byte(peer.Config),
		},
	)

	if _, err := b.app.API().Send(doc); err != nil {
		log.Printf("[tg] send file FAILED chatID=%d err=%v", chatID, err)
		return
	}

	log.Printf("[tg] sendConfig done chatID=%d duration=%s", chatID, time.Since(start))
}
