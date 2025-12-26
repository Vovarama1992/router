package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func mainKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Получить конфиг",
				"get_config",
			),
		),
	)
}
