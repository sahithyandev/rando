package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func StartTelegramPABot() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	
	if botToken == "" {
		panic("TELEGRAM_BOT_TOKEN environment variable is not set")
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}
	
	b.RegisterHandler(bot.HandlerTypeMessageText, "thought", bot.MatchTypeCommand, shareThoughtHandler)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hi! I am Sahithyan's PA.",
	})
}

func shareThoughtHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	app := "publish-thought"

    cmd := exec.Command(app, update.Message.Text)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   string(stdout),
	})
}
