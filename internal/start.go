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

	welcomeText := `🧠 <b>Привет! Я твой личный AI Cloud, умная база знаний и заметок.</b>

Скидывай сюда всё, что не хочешь потерять: мысли, статьи, документы, фото или голосовые. Я помогу найти любую запись даже по смутным воспоминаниям!

<blockquote>📥 <b>Как сохранить</b>
Начни сообщение со знака <code>+</code> (или добавь его к подписи к фото/документу):
• <code>+ Купить билеты на самолёт на пятницу</code>
• <code>+ Идея для нового пет-проекта...</code></blockquote>

<blockquote>🔍 <b>Как найти (умный поиск)</b>
Начни запрос со знака <code>?</code> - не обязательно помнить точные слова, нейросеть понимает смысл:
• <code>? что-то про поездку и билеты</code>
• <code>? рецепт пасты с грибами</code></blockquote>

<blockquote>🎙 <b>Голосовые заметки</b>
Просто запиши голосовое - я переведу речь в человеческий текст и бережно сохраню.</blockquote>

<blockquote>⏰ <b>Напомнить вовремя</b>
Напиши или скажи <code>напомни завтра в 18:00 созвониться</code> или ответь (reply) со временем на любую сохранённую запись.</blockquote>

<blockquote>🔒 <b>Приватность и безопасность</b>
Медиафайлы не скачиваются на сервер и остаются в защищённом хранилище Telegram. Но соблюдай цифровую гигиену: не отправляй в бота пароли, данные банковских карт и паспорта.</blockquote>

💡 <i>Лайфхак: если поиск найдёт ровно одну запись, я сразу отвечу реплаем на неё - ты сможешь в один клик перенестись к оригиналу в чате!</i>`

	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(message.Chat.ChatID(), welcomeText).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard),
	)
	return err
}
