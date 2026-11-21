package internal

import (
	"AICloud/internal/database"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleDeleteCallback обрабатывает нажатие на "Удалить"
func HandleDeleteCallback(ctx *th.Context, query telego.CallbackQuery) error {
	idStr := strings.TrimPrefix(query.Data, "delete:")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	userID := query.From.ID

	var messageId int
	err = database.ConnPSQL.QueryRow(ctx, "SELECT tg_message_id FROM items WHERE id = $1 AND user_id = $2", itemID, userID).Scan(&messageId)
	if err != nil {
		_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Ошибка при удалении! Не нашли заметку").WithShowAlert())
		return err
	}

	_, err = database.ConnPSQL.Exec(ctx, "DELETE FROM items WHERE id = $1 AND user_id = $2", itemID, userID)
	if err != nil {
		_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Ошибка при удалении!").WithShowAlert())
		return err
	}

	_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Заметка удалена"))

	if query.Message != nil {
		msg := query.Message.(*telego.Message)

		err = ctx.Bot().DeleteMessage(ctx, tu.Delete(msg.Chat.ChatID(), messageId))
		if err != nil {
			_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Ошибка при удалении! Не удалось удалить сообщение в чате").WithShowAlert())
			return err
		}

		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			msg.Chat.ChatID(),
			msg.MessageID,
			fmt.Sprintf("🗑 <b>Заметка #%d удалена.</b>", itemID),
		).WithParseMode(telego.ModeHTML))
		time.Sleep(2 * time.Second)
		ctx.Bot().DeleteMessage(ctx, tu.Delete(msg.Chat.ChatID(), msg.MessageID))
	}
	return nil
}

// HandleHideCallback удаляет сообщение бота при нажатии на "Скрыть"
func HandleHideCallback(ctx *th.Context, query telego.CallbackQuery) error {
	_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))

	if query.Message != nil {
		if msg, ok := query.Message.(*telego.Message); ok {
			_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(msg.Chat.ChatID(), msg.MessageID))
		}
	}
	return nil
}
