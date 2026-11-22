package internal

import (
	"AICloud/internal/database"
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleRemindCallback(ctx *th.Context, query telego.CallbackQuery) error {
	_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))

	idStr := strings.TrimPrefix(query.Data, "remind:")
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	if query.Message != nil {
		if msg, ok := query.Message.(*telego.Message); ok {
			_, err = ctx.Bot().SendMessage(ctx, tu.Message(
				msg.Chat.ChatID(),
				fmt.Sprintf("⏰ <b>Введите время для напоминания по заметке #%d:</b>\n<i>(например: 15m, 18:30, завтра, послезавтра в 14:00, 23 марта или 20 числа)</i>", itemID),
			).WithParseMode(telego.ModeHTML).WithReplyMarkup(tu.ForceReply()))
			return err
		}
	}
	return nil
}

func HandleRemindReply(ctx *th.Context, message telego.Message, rawText string, itemID int64) error {
	_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), message.MessageID))
	if message.ReplyToMessage != nil {
		_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), message.ReplyToMessage.MessageID))
	}

	duration, targetTime, err := ParseRemindTime(rawText)
	if err != nil {
		infoMsg, _ := ctx.Bot().SendMessage(ctx, tu.Message(
			message.Chat.ChatID(),
			"❌ Не понял формат времени. Напиши, например: <code>15m</code>, <code>2h</code>, <code>1d</code> или <code>18:30</code>",
		).WithParseMode(telego.ModeHTML))
		time.Sleep(4 * time.Second)
		_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(infoMsg.Chat.ChatID(), infoMsg.MessageID))
		return nil
	}

	var comment string
	_ = database.ConnPSQL.QueryRow(ctx, "SELECT comment FROM items WHERE id = $1 AND user_id = $2", itemID, message.From.ID).Scan(&comment)
	if comment == "" {
		comment = fmt.Sprintf("Заметка #%d", itemID)
	}

	bot := ctx.Bot()
	chatID := message.Chat.ChatID()

	time.AfterFunc(duration, func() {
		preview := comment
		if len([]rune(preview)) > 300 {
			preview = string([]rune(preview)[:300]) + "..."
		}

		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("🗑 Удалить").WithCallbackData(fmt.Sprintf("delete:%d", itemID)),
				tu.InlineKeyboardButton("❌ Скрыть").WithCallbackData("hide"),
			),
		)

		_, _ = bot.SendMessage(context.Background(), tu.Message(
			chatID,
			fmt.Sprintf("⏰ <b>Напоминание по заметке #%d!</b>\n\n<blockquote>%s</blockquote>", itemID, html.EscapeString(preview)),
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard))
	})

	now := time.Now()
	var timeStr string
	if targetTime.Day() == now.Day() && targetTime.Month() == now.Month() && targetTime.Year() == now.Year() {
		timeStr = "сегодня в " + targetTime.Format("15:04")
	} else if targetTime.Day() == now.AddDate(0, 0, 1).Day() && targetTime.Month() == now.AddDate(0, 0, 1).Month() && targetTime.Year() == now.AddDate(0, 0, 1).Year() {
		timeStr = "завтра в " + targetTime.Format("15:04")
	} else if targetTime.Year() == now.Year() {
		timeStr = targetTime.Format("02.01 в 15:04")
	} else {
		timeStr = targetTime.Format("02.01.2006 в 15:04")
	}

	infoMsg, _ := ctx.Bot().SendMessage(ctx, tu.Message(
		message.Chat.ChatID(),
		fmt.Sprintf("⏰ Напоминание для заметки #%d установлено на <b>%s</b>", itemID, timeStr),
	).WithParseMode(telego.ModeHTML))
	time.Sleep(3 * time.Second)
	_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(infoMsg.Chat.ChatID(), infoMsg.MessageID))

	return nil
}
