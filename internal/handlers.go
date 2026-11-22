package internal

import (
	"AICloud/internal/database"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// NewHandler регистрирует обработчики бота
func NewHandler(bh *th.BotHandler, bot *telego.Bot) {
	// Команда /start
	bh.HandleMessage(HandleStart, th.CommandEqual("start"))

	bh.HandleCallbackQuery(HandleDeleteCallback, th.CallbackDataPrefix("delete:"))
	bh.HandleCallbackQuery(HandleHideCallback, th.CallbackDataEqual("hide"))
	bh.HandleCallbackQuery(HandleFolderCallback, th.CallbackDataPrefix("folder:"))
	bh.HandleCallbackQuery(HandleRemindCallback, th.CallbackDataPrefix("remind:"))
	bh.HandleCallbackQuery(HandleProjectsCallback, th.CallbackDataEqual("projects"))

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

		if message.ReplyToMessage != nil && message.ReplyToMessage.From != nil && message.ReplyToMessage.From.IsBot {
			var folderItemID int64
			_, errFolder := fmt.Sscanf(message.ReplyToMessage.Text, "📁 Введите название папки для заметки #%d:", &folderItemID)
			if errFolder == nil && folderItemID > 0 {
				ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), message.MessageID))
				ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), message.ReplyToMessage.MessageID))

				folderName := strings.TrimSpace(rawText)
				_, err := database.ConnPSQL.Exec(ctx,
					"UPDATE items SET group_name = $1 WHERE id = $2 AND user_id = $3",
					folderName, folderItemID, message.From.ID,
				)
				if err != nil {
					_, _ = ctx.Bot().SendMessage(ctx, tu.Message(message.Chat.ChatID(), "❌ Ошибка при сохранении папки"))
					return err
				}
				infoMessage, _ := ctx.Bot().SendMessage(ctx, tu.Message(
					message.Chat.ChatID(),
					fmt.Sprintf("📁 Заметка #%d сохранена в папку <b>%s</b>", folderItemID, html.EscapeString(folderName)),
				).WithParseMode(telego.ModeHTML))
				time.Sleep(3 * time.Second)
				ctx.Bot().DeleteMessage(ctx, tu.Delete(infoMessage.Chat.ChatID(), infoMessage.MessageID))
				return nil
			}

			var remindItemID int64
			_, errRemind := fmt.Sscanf(message.ReplyToMessage.Text, "⏰ Введите время для напоминания по заметке #%d:", &remindItemID)
			if errRemind == nil && remindItemID > 0 {
				return HandleRemindReply(ctx, message, rawText, remindItemID)
			}

			return HandleEditComment(ctx, message, rawText)
		} else if message.ReplyToMessage != nil {
			return HandleEditComment(ctx, message, rawText)
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "?"); ok {
			return SearchItem(ctx, message, strings.TrimSpace(text))
		} else {
			return HandleSaveItem(ctx, message, strings.TrimSpace(rawText))
		}

	})
}
