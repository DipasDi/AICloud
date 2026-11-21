package internal

import (
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// NewHandler регистрирует обработчики бота
func NewHandler(bh *th.BotHandler, bot *telego.Bot) {
	// Команда /start
	bh.HandleMessage(HandleStart, th.CommandEqual("start"))

	// Обработка текстовых и медиа заметок (начинается с "+")
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		rawText := message.Text
		if rawText == "" {
			rawText = message.Caption
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "+"); ok {
			return HandleSaveItem(ctx, message, strings.TrimSpace(text))
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "?"); ok {
			return SearchItem(ctx, message, strings.TrimSpace(text))
		}

		return nil
	})
}
