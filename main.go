package main

import (
	"AICloud/internal"
	"AICloud/internal/database"
	"AICloud/internal/web"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	botToken := os.Getenv("TELEGRAM_TOKEN")

	go database.InitDB()
	fmt.Println("База данных запущенна")

	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)

	go internal.NewHandler(bh, bot)
	go web.StartWebServer("")
	// go ai.AnalyzerText("Блин, не понимаю, мне вроде как завтра нужно будет зайти в деконат сдать документы")

	// Stop handling updates
	defer func() { _ = bh.Stop() }()

	// Start handling updates
	_ = bh.Start()
}
