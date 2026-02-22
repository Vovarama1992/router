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
	if update.Message != nil {
		b.onMessage(update.Message)
	}
}

func (b *Bot) onMessage(msg *tgbotapi.Message) {
	log.Printf("[tg] chat=%d text=%q", msg.Chat.ID, msg.Text)

	if msg.Text == "/start" {
		m := tgbotapi.NewMessage(
			msg.Chat.ID,
			"Нажми кнопку, чтобы получить конфиг",
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

	log.Printf("[tg] CreatePeer OK chatID=%d link=%s", chatID, peer.Link)

	msg := tgbotapi.NewMessage(
		chatID,
		"Импортируй ссылку в клиент:\n\n"+peer.Link,
	)

	if _, err := b.app.API().Send(msg); err != nil {
		log.Printf("[tg] send message FAILED chatID=%d err=%v", chatID, err)
		return
	}

	log.Printf("[tg] sendConfig done chatID=%d duration=%s", chatID, time.Since(start))
}
