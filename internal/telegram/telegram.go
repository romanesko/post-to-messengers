package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"os"
	"strings"
	"xmpp/internal/common"
)

var ctx context.Context
var b *bot.Bot

func Init(context context.Context, token string) {

	ctx = context

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	var err error
	b, err = bot.New(token, opts...)

	if err != nil {
		panic(err)
	}

	go func() {
		b.Start(ctx)
		fmt.Println("Received interrupt signal, closing Telegram client...")
		os.Exit(0)
	}()

	log.Printf("Telegram bot started")

}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	message := update.Message

	if !strings.HasPrefix(message.Text, "/start") {
		return
	}

	response := common.GetEnv("TELEGRAM_WELCOME_MESSAGE", "Ваш id: <code>{{tg_chat_id}}</code>")

	if strings.Contains(response, "{{tg_chat_id}}") {
		response = strings.Replace(response, "{{tg_chat_id}}", fmt.Sprintf("%d", chatID), -1)
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      strings.ReplaceAll(response, "\\n", "\n"),
		ParseMode: models.ParseModeHTML,
	})

	if err != nil {
		log.Println(err)
	}

}

func SendMessage(chatID int64, message string) error {

	if b == nil {
		return errors.New("Teletram client is not initialized")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)

		return err
	}
	return nil
}
