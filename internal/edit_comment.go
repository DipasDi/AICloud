package internal

import (
	"AICloud/internal/ai"
	"AICloud/internal/database"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleEditComment дополняет комментарий к существующей заметке и обновляет её эмбеддинг
func HandleEditComment(ctx *th.Context, message telego.Message, newText string) error {
	return HandleEditCommentWithStatus(ctx, message, newText, nil)
}

// HandleEditCommentWithStatus дополняет комментарий, переиспользуя существующее статус-сообщение
func HandleEditCommentWithStatus(ctx *th.Context, message telego.Message, newText string, statusMsg *telego.Message) error {
	var userID int64
	if message.From != nil {
		userID = message.From.ID
	}

	if statusMsg == nil {
		var err error
		statusMsg, err = ctx.Bot().SendMessage(ctx, tu.Message(
			message.Chat.ChatID(),
			"⏳ <i>Ищу заметку для обновления...</i>",
		).WithParseMode(telego.ModeHTML))
		if err != nil {
			return fmt.Errorf("ошибка отправки статуса: %v", err)
		}
	} else {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"⏳ <i>Ищу заметку для обновления...</i>",
		).WithParseMode(telego.ModeHTML))
	}

	var (
		itemID     int64
		oldComment string
	)

	if message.ReplyToMessage != nil {
		// Если это ответ (reply) на сохранённое сообщение
		err := database.ConnPSQL.QueryRow(ctx,
			"SELECT id, comment FROM items WHERE user_id = $1 AND tg_message_id = $2",
			userID, message.ReplyToMessage.MessageID,
		).Scan(&itemID, &oldComment)
		if err != nil {
			_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
				statusMsg.Chat.ChatID(),
				statusMsg.MessageID,
				"❌ <b>Заметка не найдена</b>\nСделай ответ (reply) на сохранённое сообщение, чтобы дополнить его.",
			).WithParseMode(telego.ModeHTML))
			return err
		}
	} else {
		// Иначе берём последнюю запись пользователя
		err := database.ConnPSQL.QueryRow(ctx,
			"SELECT id, comment FROM items WHERE user_id = $1 ORDER BY id DESC LIMIT 1",
			userID,
		).Scan(&itemID, &oldComment)
		if err != nil {
			_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
				statusMsg.Chat.ChatID(),
				statusMsg.MessageID,
				"❌ <b>Заметка не найдена</b>\nСделай ответ (reply) на сохранённое сообщение, чтобы дополнить его.",
			).WithParseMode(telego.ModeHTML))
			return err
		}
	}

	cleanNewText := strings.TrimSpace(newText)
	newComment := strings.TrimSpace(oldComment + " " + cleanNewText)

	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		"⏳ <i>Получаю новый эмбеддинг от ИИ...</i>",
	).WithParseMode(telego.ModeHTML))

	embedding, err := ai.GetEmbedding(newComment)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка ИИ при получении эмбеддинга</b>",
		).WithParseMode(telego.ModeHTML))
		return fmt.Errorf("ошибка получения Embedding от ИИ: %v", err)
	}

	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		"💾 <i>Сохраняю изменения в базу данных...</i>",
	).WithParseMode(telego.ModeHTML))

	_, err = database.ConnPSQL.Exec(ctx,
		"UPDATE items SET comment = $1, embedding = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4",
		newComment, FormatVector(embedding), itemID, userID,
	)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка сохранения в базу данных</b>",
		).WithParseMode(telego.ModeHTML))
		log.Println("Ошибка обновления в базе: ", err)
		return fmt.Errorf("ошибка обновления записи: %v", err)
	}

	previewText := newComment
	if len([]rune(previewText)) > 120 {
		previewText = string([]rune(previewText)[:120]) + "..."
	}

	finalText := fmt.Sprintf(
		"✅ <b>Заметка #%d обновлена!</b>\n\n"+
			"<blockquote>%s</blockquote>\n"+
			"💡 <i>Новый текст и вектор успешно сохранены.</i>",
		itemID,
		html.EscapeString(previewText),
	)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⏰ Напомнить").WithCallbackData(fmt.Sprintf("remind:%d", itemID)),
			tu.InlineKeyboardButton("📁 В папку").WithCallbackData(fmt.Sprintf("folder:%d", itemID)),
			tu.InlineKeyboardButton("🗑 Удалить").WithCallbackData(fmt.Sprintf("delete:%d", itemID)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("❌ Скрыть").WithCallbackData("hide"),
		),
	)

	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		finalText,
	).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard))

	return nil
}
