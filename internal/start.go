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
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🚀 Другие бесплатные проекты").WithCallbackData("projects"),
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

	welcomeText := `🧠 <b>Привет! Я твой личный AI Cloud - умная база знаний и персональное облако в Telegram.</b>

Скидывай сюда всё, что важно: мысли, ссылки, статьи, файлы, код, фото или голосовые. Я нахожу нужные записи даже по смутным воспоминаниям!

<blockquote>📥 <b>Как сохранять</b>
• Начни сообщение со знака <code>+</code>:
  <code>+ Купить билеты в Казань на пятницу</code>
• Или пришли файл, фото, документ или ссылку с подписью через <code>+</code>.
• Любое голосовое я сам переведу в текст и сохраню как заметку.</blockquote>

<blockquote>🔍 <b>Как искать (семантический поиск)</b>
Начни запрос со знака <code>?</code> - точные слова помнить не нужно, поиск понимает смысл:
• <code>? что-то про поездку и билеты</code>
• <code>? конфиг базы данных postgres</code></blockquote>

<blockquote>📁 <b>Папки и категории</b>
Жми кнопку <b>«📁 В папку»</b> под любой заметкой и отправляй название (Работа, Учеба, Код, Рецепты). Заметки сразу сортируются по категориям.</blockquote>

<blockquote>⏰ <b>Умные напоминания</b>
Жми кнопку <b>«⏰ Напомнить»</b> под заметкой. Я понимаю человеческий язык и любые форматы:
• Интервалы: <code>15m</code>, <code>2h</code>, <code>через 3 дня</code>
• Точное время: <code>18:30</code>, <code>в 19:00</code>
• Дни: <code>завтра</code>, <code>завтра в 15:30</code>, <code>послезавтра</code>
• Календарь: <code>23 марта в 19:00</code>, <code>20 числа</code>, <code>23.10</code></blockquote>

<blockquote>✏️ <b>Быстрое дополнение заметок</b>
Ответь (reply) на своё сообщение или заметку, чтобы дописать новый текст. Я обновлю запись и пересчитаю векторный поиск.</blockquote>

<blockquote>🔒 <b>Приватность и безопасность</b>
Медиафайлы не скачиваются на сторонние серверы и остаются в хранилище Telegram. Но соблюдай цифровую гигиену: не храни здесь пароли и данные банковских карт.</blockquote>

💡 <i>Лайфхак: если поиск найдет ровно одну запись, я сразу отвечу реплаем на оригинал - ты сможешь в один клик прыгнуть к сообщению в чате!</i>`

	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(message.Chat.ChatID(), welcomeText).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboard),
	)
	return err
}
