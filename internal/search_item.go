package internal

import (
	"AICloud/internal/ai"
	"AICloud/internal/database"
	"AICloud/models"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SearchItem(ctx *th.Context, message telego.Message, text string) error {

	ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), message.MessageID))

	// Отправляем начальный статус пользователю
	statusMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		message.Chat.ChatID(),
		"⏳ <i>Ищу...</i>",
	).WithParseMode(telego.ModeHTML))
	if err != nil {
		return fmt.Errorf("ошибка отправки статуса: %v", err)
	}

	var userID int64
	if message.From != nil {
		userID = message.From.ID
	}

	// Сначала ищем по тексту без ИИ (полнотекстовый поиск PostgreSQL)
	textSearchQuery := `
		SELECT 
			id,
			target_type,
			tg_file_id,
			tg_message_id,
			file_name,
			comment,
			group_name,
			created_at,
			ts_rank(search_vector, plainto_tsquery('russian', $2)) AS rank
		FROM items
		WHERE user_id = $1
		  AND search_vector @@ plainto_tsquery('russian', $2)
		ORDER BY rank DESC
		LIMIT 5;
	`

	rows, err := database.ConnPSQL.Query(ctx, textSearchQuery, userID, text)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var it models.Item
		err := rows.Scan(
			&it.ID,
			&it.TargetType,
			&it.TgFileID,
			&it.TgMessageID,
			&it.FileName,
			&it.Comment,
			&it.GroupName,
			&it.CreatedAt,
			&it.Rank,
		)
		if err != nil {
			return err
		}
		items = append(items, it)
	}

	// Если нашли по тексту выводим результат и завершаем
	if len(items) > 0 {
		sendSearchResults(ctx, statusMsg, "🔍 <b>Найдено по тексту:</b>", items)
		return nil
	}

	// Если по тексту ничего не нашлось ищем по смыслу через ИИ (pgvector)
	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		statusMsg.Chat.ChatID(),
		statusMsg.MessageID,
		"🔍 В тексте совпадений нет.\n🧠 <i>Ищу по смыслу через нейросеть...</i>",
	).WithParseMode(telego.ModeHTML))

	embedding, err := ai.GetEmbedding(text)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка ИИ при получении эмбеддинга</b>",
		).WithParseMode(telego.ModeHTML))
		return fmt.Errorf("ошибка получения Embedding от ИИ: %v", err)
	}

	vectorSearchQuery := `
		SELECT 
			id,
			target_type,
			tg_file_id,
			tg_message_id,
			file_name,
			comment,
			group_name,
			created_at,
			1 - (embedding <=> $2) AS rank
		FROM items
		WHERE user_id = $1
		  AND embedding IS NOT NULL
		ORDER BY embedding <=> $2
		LIMIT 5;
	`

	vRows, err := database.ConnPSQL.Query(ctx, vectorSearchQuery, userID, FormatVector(embedding))
	if err != nil {
		return fmt.Errorf("ошибка векторного поиска: %v", err)
	}
	defer vRows.Close()

	var aiItems []models.Item
	for vRows.Next() {
		var it models.Item
		err := vRows.Scan(
			&it.ID,
			&it.TargetType,
			&it.TgFileID,
			&it.TgMessageID,
			&it.FileName,
			&it.Comment,
			&it.GroupName,
			&it.CreatedAt,
			&it.Rank,
		)
		if err != nil {
			return err
		}

		if it.Rank >= 0.25 {
			aiItems = append(aiItems, it)
		}
	}

	hideKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("❌ Скрыть").WithCallbackData("hide"),
		),
	)

	if len(aiItems) == 0 {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			fmt.Sprintf("🔍 По запросу «%s» ничего не найдено ни в тексте, ни по смыслу.", html.EscapeString(text)),
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
		time.Sleep(4 * time.Second)
		ctx.Bot().DeleteMessage(ctx, tu.Delete(message.Chat.ChatID(), statusMsg.MessageID))

		return nil
	}

	sendSearchResults(ctx, statusMsg, "🧠 <b>Найдено по смыслу (ИИ):</b>", aiItems)
	return nil
}

// sendSearchResults отправляет результаты: реплаем на исходное сообщение
func sendSearchResults(ctx *th.Context, statusMsg *telego.Message, title string, items []models.Item) {
	chatID := statusMsg.Chat.ChatID()
	hideKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("❌ Скрыть").WithCallbackData("hide"),
		),
	)

	if len(items) == 1 {
		it := items[0]
		_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(chatID, statusMsg.MessageID))

		msgText := fmt.Sprintf("%s\n\n%s", title, formatSingleItem(it))
		msg := tu.Message(chatID, msgText).
			WithParseMode(telego.ModeHTML).
			WithReplyMarkup(hideKeyboard)

		if it.TgMessageID != 0 {
			msg.WithReplyParameters(&telego.ReplyParameters{
				MessageID:                int(it.TgMessageID),
				AllowSendingWithoutReply: true,
			})
		}

		_, _ = ctx.Bot().SendMessage(ctx, msg)
		return
	}

	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		chatID,
		statusMsg.MessageID,
		formatSearchResults(title, items),
	).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
}

// formatSingleItem форматирует одну найденную запись
func formatSingleItem(it models.Item) string {
	icon := "💬"
	switch it.TargetType {
	case "file":
		icon = "📄"
	case "photo":
		icon = "🖼"
	case "voice":
		icon = "🎙"
	case "link":
		icon = "🔗"
	}

	header := ""
	if it.FileName != "" {
		header = fmt.Sprintf("<b>%s</b> ", html.EscapeString(it.FileName))
	}

	comment := it.Comment
	if len([]rune(comment)) > 200 {
		comment = string([]rune(comment)[:200]) + "..."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s<code>#%d</code>\n", icon, header, it.ID))
	if comment != "" {
		sb.WriteString(fmt.Sprintf("<blockquote>%s</blockquote>\n", html.EscapeString(comment)))
	}
	sb.WriteString(fmt.Sprintf("📅 <i>%s</i>", it.CreatedAt.Format("02.01 в 15:04")))

	return sb.String()
}

// formatSearchResults форматирует список найденных записей в красивый HTML
func formatSearchResults(title string, items []models.Item) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s <i>(%d)</i>\n\n", title, len(items)))

	for i, it := range items {
		icon := "💬"
		switch it.TargetType {
		case "file":
			icon = "📄"
		case "photo":
			icon = "🖼"
		case "voice":
			icon = "🎙"
		case "link":
			icon = "🔗"
		}

		header := ""
		if it.FileName != "" {
			header = fmt.Sprintf("<b>%s</b> ", html.EscapeString(it.FileName))
		}

		comment := it.Comment
		if len([]rune(comment)) > 300 {
			comment = string([]rune(comment)[:300]) + "..."
		}

		sb.WriteString(fmt.Sprintf("%d. %s %s<code>#%d</code>\n", i+1, icon, header, it.ID))
		if comment != "" {
			sb.WriteString(fmt.Sprintf("<blockquote>%s</blockquote>\n", html.EscapeString(comment)))
		}
		sb.WriteString(fmt.Sprintf("📅 <i>%s</i>\n\n", it.CreatedAt.Format("02.01 в 15:04")))
	}

	return strings.TrimSpace(sb.String())
}
