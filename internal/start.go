package internal

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleStart обрабатывает команду /start
func HandleStart(ctx *th.Context, message telego.Message) error {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Поддержка").WithURL("https://t.me/dipas_di"),
			tu.InlineKeyboardButton("ТГК разработчика").WithURL("https://t.me/DipasDEV"),
		),
	)

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
}
