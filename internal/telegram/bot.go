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

	peer, err := b.svc.CreatePeer(context.Background(), chatID)
	if err != nil {
		b.app.API().Send(
			tgbotapi.NewMessage(chatID, "Ошибка создания конфига"),
		)
		return
	}

	instruction := `Сейчас отправлю ссылку отдельным сообщением.
Как подключиться:
Android — установите Hiddify из Google Play.
iPhone — установите Streisand из App Store.
Компьютер — установите Hiddify: https://hiddify.org/en/download-hiddify/
После установки:
Скопируйте ссылку
Откройте приложение
Нажмите «Создать подключение»
Выберите «Импорт из буфера обмена»`

	b.app.API().Send(tgbotapi.NewMessage(chatID, instruction))
	b.app.API().Send(tgbotapi.NewMessage(chatID, peer.Link))

	log.Printf("[tg] sendConfig done chatID=%d duration=%s", chatID, time.Since(start))
}
