package internal

import (
	"AICloud/internal/ai"
	"AICloud/internal/database"
	"AICloud/models"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Обработчик команд
func NewHandler(bh *th.BotHandler, bot *telego.Bot) {
	// Собираем кнопки start в один ряд
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Поддержка").WithURL("https://t.me/dipas_di"),
			tu.InlineKeyboardButton("ТГК разработчика").WithURL("https://t.me/DipasDEV"),
		),
	)

	//Комманда Start
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		_, _, args := tu.ParseCommand(message.Text)

		if len(args) > 0 {

			payload := args[0]
			_, err := ctx.Bot().SendMessage(ctx, tu.Messagef(
				message.Chat.ChatID(),
				"Запуск с параметром: <code>%s</code>", payload,
			).WithParseMode(telego.ModeHTML))
			return err

		}

		welcomeText := `<b>Привет! Я твоя личная база знаний и напоминалка.</b>

Скидывай всё сюда - я сам расшифрую голосовые, вытащу текст со скриншотов и разложу по полкам.

<blockquote><b>Сохранить</b>
Перешли пост, кинь ссылку, фото, PDF, голосовое или просто напиши длинную мысль (либо начни с <code>+</code>).</blockquote>

<blockquote><b>Найти по смыслу</b>
Не нужно помнить точные слова:
• <code>найди рецепт пасты</code>
• <code>? та статья про настройку docker</code>
• <code>/find билеты</code></blockquote>

<blockquote><b>Напомнить вовремя</b>
Напиши или скажи <code>напомни завтра в 18:00 созвониться</code> или ответь (reply) со временем на любую сохранённую запись.</blockquote>

<i>Лайфхак: если сделать reply на сохранённый пост обычным текстом, я прикреплю его как твою заметку.</i>`

		_, err := ctx.Bot().SendMessage(
			ctx,
			tu.Message(message.Chat.ChatID(), welcomeText).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard),
		)
		return err
	}, th.CommandEqual("start"))

	// Обработка команда, медиа или текст
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		rawText := message.Text
		if rawText == "" {
			rawText = message.Caption
		}

		if text, ok := strings.CutPrefix(strings.TrimSpace(rawText), "+"); ok {
			return HandleSaveItem(ctx, message, strings.TrimSpace(text))
		}

		return nil
	})
}

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

// Сохранение
func HandleSaveItem(ctx *th.Context, message telego.Message, text string) error {
	var err error

	//Извлечения данных медиа и получение типа
	mediaInfo := ExtractMediaInfo(message)

	//Отправляем сообщение пользователю
	statusMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		message.Chat.ChatID(),
		"⏳ <i>Получаю эмбеддинг от ИИ...</i>",
	).WithParseMode(telego.ModeHTML))
	if err != nil {
		return fmt.Errorf("ошибка отправки статуса: %v", err)
	}

	//Отправляем и получаем вектор от ии
	embedding, err := ai.GetEmbedding(text)
	if err != nil {
		_, _ = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			statusMsg.Chat.ChatID(),
			statusMsg.MessageID,
			"❌ <b>Ошибка ИИ при получении эмбеддинга</b>",
		).WithParseMode(telego.ModeHTML))
		return fmt.Errorf("ошибка получения Embedding от ИИ: %v", err)
	}

	//Сообщем что мы сохраняем в базу данных
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

func FormatVector(v []float32) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, val := range v {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 32))
	}
	sb.WriteByte(']')
	return sb.String()
}

func GetFileBytes(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return body
}

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
	var analiz bool
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
		analiz = true
	}

	_, err = bot.EditMessageText(ctx, tu.EditMessageText(chatID, statusMsg.MessageID, finalMsg).WithParseMode(telego.ModeHTML))
	if err != nil {
		log.Println("Ошибка редактирования сообщения:", err)
	}

	if analiz {

	}
}
