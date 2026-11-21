package internal

import (
	"AICloud/internal/ai"
	"AICloud/internal/database"
	"fmt"
	"html"
	"log"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleSaveItem сохраняет входящую заметку (текст или с медиа) в базу данных
func HandleSaveItem(ctx *th.Context, message telego.Message, text string) error {
	return HandleSaveItemWithStatus(ctx, message, text, nil)
}

// HandleSaveItemWithStatus сохраняет заметку, переиспользуя существующее статус-сообщение (если есть)
func HandleSaveItemWithStatus(ctx *th.Context, message telego.Message, text string, statusMsg *telego.Message) error {
	var err error

	// Извлечение данных медиа и получение типа
	mediaInfo := ExtractMediaInfo(message)

	// Отправляем или обновляем начальный статус пользователю
	if statusMsg == nil {
		statusMsg, err = ctx.Bot().SendMessage(ctx, tu.Message(
			message.Chat.ChatID(),
			"⏳ <i>Получаю эмбеддинг от ИИ...</i>",
		).WithParseMode(telego.ModeHTML))
		if err != nil {
			return fmt.Errorf("ошибка отправки статуса: %v", err)
		}
	} else {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"⏳ <i>Получаю эмбеддинг от ИИ...</i>",
		).WithParseMode(telego.ModeHTML))
	}

	// Отправляем и получаем вектор от ИИ
	embedding, err := ai.GetEmbedding(text)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка ИИ при получении эмбеддинга</b>",
		).WithParseMode(telego.ModeHTML))
		return fmt.Errorf("ошибка получения Embedding от ИИ: %v", err)
	}

	// Сообщаем о сохранении в базу данных
	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		"💾 <i>Сохраняю в базу данных...</i>",
	).WithParseMode(telego.ModeHTML))

	args := make([]any, 0, 10)
	var query string

	var userID int64
	if message.From != nil {
		userID = message.From.ID
	}

	query = `INSERT INTO items (
				user_id,
				chat_id,
				target_type,
				tg_file_id,
				tg_message_id,
				file_name,
				comment,
				group_name,
				embedding,
				embedding_model
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	args = append(args,
		userID,                   // $1: user_id
		message.Chat.ID,          // $2: chat_id
		mediaInfo.TargetType,     // $3: target_type
		mediaInfo.FileID,         // $4: tg_file_id
		message.MessageID,        // $5: tg_message_id
		mediaInfo.FileName,       // $6: file_name
		text,                     // $7: comment
		"",                       // $8: group_name
		FormatVector(embedding),  // $9: embedding
		os.Getenv("EMBED_MODEL"), // $10: embedding_model
	)

	var itemID int64
	err = database.ConnPSQL.QueryRow(ctx, query, args...).Scan(&itemID)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка сохранения в базу данных</b>",
		).WithParseMode(telego.ModeHTML))
		log.Println("Ошибка записи в базу: ", err)
		return fmt.Errorf("ошибка записи: %v", err)
	}

	previewText := text
	if len([]rune(previewText)) > 120 {
		previewText = string([]rune(previewText)[:120]) + "..."
	}

	finalText := fmt.Sprintf(
		"✅ <b>Заметка #%d сохранена!</b>\n\n"+
			"<blockquote>%s</blockquote>\n"+
			"💡 <i>Сможешь найти её по смыслу через поиск.</i>",
		itemID,
		html.EscapeString(previewText),
	)

	if mediaInfo.TargetType == "voice" {
		finalText = fmt.Sprintf(
			"🎙 <b>Перевёл с голосового на человеческий (#%d):</b>\n\n"+
				"<blockquote>%s</blockquote>\n"+
				"<i>Если что-то переврал - виновата дикция, но я старался.</i>",
			itemID,
			html.EscapeString(previewText),
		)
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⏰ Напомнить").WithCallbackData(fmt.Sprintf("remind:%d", itemID)),
			tu.InlineKeyboardButton("📁 В папку").WithCallbackData(fmt.Sprintf("folder:%d", itemID)),
			tu.InlineKeyboardButton("🗑 Удалить").WithCallbackData(fmt.Sprintf("delete:%d", itemID)),
		),
	)

	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		finalText,
	).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard))

	return err
}
