package internal

import (
	"AICloud/internal/ai"
	"AICloud/models"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// ExtractMediaInfo извлекает тип, ID файла и имя из сообщения Telegram
func ExtractMediaInfo(msg telego.Message) models.MediaInfo {
	switch {
	case len(msg.Photo) > 0:
		return models.MediaInfo{
			TargetType: "photo",
			FileID:     msg.Photo[len(msg.Photo)-1].FileID,
		}
	case msg.Document != nil:
		return models.MediaInfo{
			TargetType: "file",
			FileID:     msg.Document.FileID,
			FileName:   msg.Document.FileName,
		}
	case msg.Voice != nil:
		return models.MediaInfo{
			TargetType: "voice",
			FileID:     msg.Voice.FileID,
		}
	default:
		return models.MediaInfo{
			TargetType: "message",
		}
	}
}

// GetFileBytes скачивает файл по ссылке и возвращает байты
func GetFileBytes(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

// FileVoice обрабатывает голосовые сообщения
func FileVoice(bh *th.BotHandler, bot *telego.Bot, ctx *th.Context, message telego.Message) {
	chatID := message.Chat.ChatID()

	_ = bot.SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))
	statusMsg, err := bot.SendMessage(ctx, tu.Message(chatID, "🎧 <i>Навострил уши, слушаю...</i>").WithParseMode(telego.ModeHTML))
	if err != nil {
		log.Println("Ошибка отправки статуса:", err)
		return
	}

	tgFile, err := bot.GetFile(ctx, &telego.GetFileParams{FileID: message.Voice.FileID})
	if err != nil {
		log.Println("Ошибка GetFile:", err)
		_, _ = bot.EditMessageText(ctx, tu.EditMessageText(chatID, statusMsg.MessageID,
			"📡 <b>Телеграм зажал файл</b>\nНе смог скачать голосовое. Попробуй переслать ещё раз.").WithParseMode(telego.ModeHTML))
		return
	}

	url := bot.FileDownloadURL(tgFile.FilePath)
	audioBytes := GetFileBytes(url)

	text, err := ai.AnalysisVoice(ctx, audioBytes)

	var finalMsg string
	cleanText := strings.TrimSpace(text)

	switch {
	case err != nil:
		log.Println("Ошибка AI:", err)
		finalMsg = "💥 <b>Шестерёнки заклинило</b>\n" +
			"Разраб опять что-то накрутил в бэкенде, и нейросеть ушла на перекур.\n\n" +
			"<i>Кинь голосовое чуть позже - я уже пинаю его, чтобы починил.</i>"

	case cleanText == "" || strings.Contains(cleanText, "Продолжение следует..."):
		finalMsg = "🙉 <b>Режим телепатии пока в бете</b>\n" +
			"Там либо тишина, либо шифровка для своих - разобрать не вышло.\n\n" +
			"<i>Попробуй сказать погромче или поближе к микрофону.</i>"

	default:
		finalMsg = fmt.Sprintf(
			"🎙 <b>Перевёл с голосового на человеческий:</b>\n"+
				"<blockquote>%s</blockquote>\n"+
				"<i>Если что-то переврал - виновата дикция, но я старался. Начинаю разбор</i>",
			html.EscapeString(cleanText),
		)
	}

	_, err = bot.EditMessageText(ctx, tu.EditMessageText(chatID, statusMsg.MessageID, finalMsg).WithParseMode(telego.ModeHTML))
	if err != nil {
		log.Println("Ошибка редактирования сообщения:", err)
	}
}
