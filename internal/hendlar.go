package internal

import (
	"AICloud/internal/ai"
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

// Обработчик команд
func NewHandler(bh *th.BotHandler, bot *telego.Bot) {

	// Собираем кнопки start в один ряд
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Поддержка").WithURL("https://t.me/dipas_di"),
			tu.InlineKeyboardButton("ТГК разработчика").WithURL("https://t.me/DipasDEV"),
		),
	)

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {

		switch {
		case message.Voice != nil:
			fmt.Println("voice")
			FileVoice(bh, bot, ctx, message)

		case message.Photo != nil:
			fmt.Println("photo")
		case message.Audio != nil:
			fmt.Println("audio")
		case message.Document != nil:
			fmt.Println("document")
		case message.Sticker != nil:
			fmt.Println("sticker")
		case message.Video != nil:
			fmt.Println("video")
		}

		return nil
	}, th.AnyMessageWithMedia())

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

	text, err := ai.AINalysisVoice(ctx, audioBytes)

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
				"<i>Если что-то переврал - виновата дикция, но я старался.</i>",
			html.EscapeString(cleanText),
		)
	}

	// 4. Плавно заменяем "Слушаю..." на готовый результат
	_, err = bot.EditMessageText(ctx, tu.EditMessageText(chatID, statusMsg.MessageID, finalMsg).WithParseMode(telego.ModeHTML))
	if err != nil {
		log.Println("Ошибка редактирования сообщения:", err)
	}
}
