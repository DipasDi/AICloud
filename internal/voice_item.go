package internal

import (
	"AICloud/internal/ai"
	"fmt"
	"log"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleVoiceMessage обрабатывает голосовые сообщения (заметки, реплаи и голосовой поиск)
func HandleVoiceMessage(ctx *th.Context, message telego.Message) error {
	chatID := message.Chat.ChatID()

	_ = ctx.Bot().SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))

	hideKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("❌ Скрыть").WithCallbackData("hide"),
		),
	)

	statusMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		chatID,
		"🎧 <i>Навострил уши, слушаю...</i>",
	).WithParseMode(telego.ModeHTML))
	if err != nil {
		log.Println("Ошибка отправки статуса:", err)
		return fmt.Errorf("ошибка отправки статуса: %w", err)
	}

	tgFile, err := ctx.Bot().GetFile(ctx, &telego.GetFileParams{FileID: message.Voice.FileID})
	if err != nil {
		log.Println("Ошибка GetFile:", err)
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			chatID,
			statusMsg.MessageID,
			"📡 <b>Телеграм зажал файл</b>\nНе смог скачать голосовое. Попробуй переслать ещё раз.",
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
		return err
	}

	url := ctx.Bot().FileDownloadURL(tgFile.FilePath)
	audioBytes, err := GetFileBytes(url)
	if err != nil {
		log.Println("Ошибка скачивания файла:", err)
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			chatID,
			statusMsg.MessageID,
			"📡 <b>Телеграм зажал файл</b>\nНе смог скачать голосовое. Попробуй переслать ещё раз.",
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
		return err
	}

	text, err := ai.AnalysisVoice(ctx, audioBytes)
	cleanText := strings.TrimSpace(text)

	switch {
	case err != nil:
		log.Println("Ошибка AI:", err)
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			chatID,
			statusMsg.MessageID,
			"💥 <b>Шестерёнки заклинило</b>\n"+
				"Разраб опять что-то накрутил в бэкенде, и нейросеть ушла на перекур.\n\n"+
				"<i>Кинь голосовое чуть позже - я уже пинаю его, чтобы починил.</i>",
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
		return err

	case cleanText == "" || strings.Contains(cleanText, "Продолжение следует..."):
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			chatID,
			statusMsg.MessageID,
			"🙉 <b>Режим телепатии пока в бете</b>\n"+
				"Там либо тишина, либо шифровка для своих - разобрать не вышло.\n\n"+
				"<i>Попробуй сказать погромче или поближе к микрофону.</i>",
		).WithParseMode(telego.ModeHTML).WithReplyMarkup(hideKeyboard))
		return nil
	}

	if message.ReplyToMessage != nil {
		return HandleEditCommentWithStatus(ctx, message, cleanText, statusMsg)
	}

	_ = ctx.Bot().SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))
	_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
		chatID,
		statusMsg.MessageID,
		"🧠 <i>Понимаю суть запроса...</i>",
	).WithParseMode(telego.ModeHTML))

	intent, err := ai.AnalyzerText(ctx, cleanText)
	if err != nil {
		log.Println("Ошибка анализа намерения, сохраняем как обычную заметку:", err)
		return HandleSaveItemWithStatus(ctx, message, cleanText, statusMsg)
	}

	switch intent.Action {
	case "search":
		searchQuery := strings.TrimSpace(intent.Query)
		if searchQuery == "" {
			searchQuery = cleanText
		}
		_ = ctx.Bot().DeleteMessage(ctx, tu.Delete(chatID, statusMsg.MessageID))
		return SearchItem(ctx, message, searchQuery)

	case "save", "remind":
		noteText := strings.TrimSpace(intent.Query)
		if noteText == "" {
			noteText = cleanText
		}
		return HandleSaveItemWithStatus(ctx, message, noteText, statusMsg)

	default:
		return HandleSaveItemWithStatus(ctx, message, cleanText, statusMsg)
	}
}
