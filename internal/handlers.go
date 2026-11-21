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

	bh.HandleCallbackQuery(HandleDeleteCallback, th.CallbackDataPrefix("delete:"))
	bh.HandleCallbackQuery(HandleHideCallback, th.CallbackDataEqual("hide"))

	// Обработка текстовых, медиа и голосовых сообщений
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		// Обработка голосовых сообщений
		if message.Voice != nil {
			return HandleVoiceMessage(ctx, message)
		}

		rawText := message.Text
		if rawText == "" {
			rawText = message.Caption
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "+"); ok {
			return HandleSaveItem(ctx, message, strings.TrimSpace(text))
		}

		if message.ReplyToMessage != nil {
			return HandleEditComment(ctx, message, rawText)
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "?"); ok {
			return SearchItem(ctx, message, strings.TrimSpace(text))
		} else {
			return HandleSaveItem(ctx, message, strings.TrimSpace(rawText))
		}

	})
}
