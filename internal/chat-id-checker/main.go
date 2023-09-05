package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send /id to the bot after the bot has been started - you will get your username and chat id

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil && update.Message.Text == "/id" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%s: %d", update.Message.Chat.Username, update.Message.Chat.ID),
		})
		if err != nil {
			log.Println("Error", err)
		}
	}

	if update.ChannelPost != nil && update.ChannelPost.Text == "/id" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChannelPost.Chat.ID,
			Text:   fmt.Sprintf("%d", update.ChannelPost.Chat.ID),
		})
		if err != nil {
			log.Println("Error", err)
		}
	}
}
