package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const OwnerUserId = int64(6328469595)

func StartTelegramPABot(ctx context.Context) {
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

	b.RegisterHandler(bot.HandlerTypeMessageText, "/thought", bot.MatchTypePrefix, shareThoughtHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if OwnerUserId == update.Message.From.ID {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Hi! I am your PA. You know me.",
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Hi! I am Sahithyan's PA. You don't have any authority over me for now.",
		})
	}
}

func shareThoughtHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if OwnerUserId != update.Message.From.ID {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You are not allowed to share thoughts. Only Sahithyan can.",
		})
		return
	}
	
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 2)
	command := parts[0]
	
	if command != "/thought" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Invalid command. Use /thought to share a thought.",
		})
		return
	}
	
	var args string
	if len(parts) > 1 {
		args = parts[1]
	}
	
	if args == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a thought to share.",
		})
		return
	}

	cmd := exec.Command("publish-thought", args)
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
func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if OwnerUserId == update.Message.From.ID {
		helpText := `Available Commands:
/start - Start interacting with the bot
/thought - Share a new thought`

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   helpText,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Currently you can't use this bot. Only Sahithyan can.",
		})
	}
}
